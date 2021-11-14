package fflogs

import (
	"context"
	stdjson "encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	graphql "github.com/hasura/go-graphql-client"
	json "github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

type Client interface {
	GetTimesFromFightAndID(ctx context.Context, code string,
		fightID int) (startTime, endTime float64, err error)
	GetEvents(ctx context.Context, code string, fightID int, startTime,
		endTime float64, dataType EventDataType, hostilityType HostilityType) ([]*RawBuffEvent, error)
}

type ClientImplementation struct {
	c        *graphql.Client
	h        *resty.Client
	clientID string
	secret   string
	token    *oauth2.Token
}

var _ Client = &ClientImplementation{}

const fflogsGraphQLEndpoint = "https://www.fflogs.com/api/v2/client"
const fflogsOauthEndpoint = "https://www.fflogs.com/oauth/token"

// TODO: caching
func NewClient(clientID, secret string, token ...string) (*ClientImplementation, error) {
	c := &ClientImplementation{
		h:        resty.New(),
		clientID: clientID,
		secret:   secret,
	}
	if len(token) == 1 {
		c.token = &oauth2.Token{AccessToken: token[0]}
		c.c = graphql.NewClient(fflogsGraphQLEndpoint, oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(c.token)))
	} else {
		c.c = graphql.NewClient(fflogsGraphQLEndpoint,
			oauth2.NewClient(context.Background(), c))
		v, err := c.getCreds()
		if err != nil {
			return nil, err
		}
		c.token = v
		time.AfterFunc(v.Expiry.Sub(time.Now())-time.Minute, c.getCredsWrapper)
	}
	return c, nil
}

func (c *ClientImplementation) getCredsWrapper() {
	v, err := c.getCreds()
	if err != nil {
		panic(fmt.Errorf("fatal: unable to refresh bearer token: %v", err))
	}
	c.token = v
	time.AfterFunc(v.Expiry.Sub(time.Now())-time.Minute, c.getCredsWrapper)
}

func (c *ClientImplementation) Token() (*oauth2.Token, error) {
	return c.token, nil
}

func (c *ClientImplementation) getCreds() (*oauth2.Token, error) {
	v := &oauth2.Token{}
	err := resty.Backoff(func() (*resty.Response, error) {
		t0 := time.Now()
		defer func() { log.Printf("Oauth2 took %.2fs\n", time.Since(t0).Seconds()) }()
		resp, err := c.h.R().
			SetBasicAuth(c.clientID, c.secret).
			SetFormData(map[string]string{"grant_type": `client_credentials`}).
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			Post(fflogsOauthEndpoint)
		if err != nil {
			return nil, err
		}
		return resp, json.Unmarshal(resp.Body(), v)
	}, resty.RetryConditions([]resty.RetryConditionFunc{func(r *resty.Response, e error) bool {
		return r.StatusCode() != http.StatusOK
	}}))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *ClientImplementation) GetTimesFromFightAndID(ctx context.Context, code string, fightID int) (startTime, endTime float64, err error) {
	t0 := time.Now()
	defer func() { log.Printf("GetTimesFromFightAndID took %.2fs\n", time.Since(t0).Seconds()) }()
	var query struct {
		ReportData struct {
			Report struct {
				Fights []struct {
					StartTime graphql.Float
					EndTime   graphql.Float
				} `graphql:"fights(fightIDs:[$fightID])"`
			} `graphql:"report(code:$code)"`
		} `graphql:"reportData"`
	}
	err = c.c.Query(ctx, &query, map[string]interface{}{"code": graphql.String(code), "fightID": graphql.Int(fightID)})
	if err != nil {
		return 0, 0, err
	}
	return float64(query.ReportData.Report.Fights[0].StartTime), float64(query.ReportData.Report.Fights[0].EndTime), nil
}

type Ability struct {
	Name string
	Icon string
}

func (c *ClientImplementation) GetAllAbilities(ctx context.Context) (map[int]Ability, error) {
	t0 := time.Now()
	defer func() { log.Printf("GetAllAbilities took %.2fs\n", time.Since(t0).Seconds()) }()
	// {
	//   gameData {
	//     abilities(limit:100, page:1) {
	//       has_more_pages
	//       data{
	//         id
	//         icon
	//         name
	//       }
	//     }
	//   }
	// }
	ret := map[int]Ability{}
	type query struct {
		GameData struct {
			Abilities struct {
				More graphql.Boolean `graphql:"has_more_pages"`
				Data []struct {
					ID   graphql.Int    `graphql:"id"`
					Icon graphql.String `graphql:"icon"`
					Name graphql.String `graphql:"name"`
				} `graphql:"data"`
			} `graphql:"abilities(limit:$limit,page:$page)"`
		} `graphql:"gameData"`
	}
	cont := true
	page := 1
	t00 := time.Now()
	for cont {
		t0 := time.Now()
		q := query{}
		err := c.c.Query(ctx, &q, map[string]interface{}{
			"limit": graphql.Int(100),
			"page":  graphql.Int(page),
		})
		if err != nil {
			return nil, err
		}
		for _, v := range q.GameData.Abilities.Data {
			ret[int(v.ID)] = Ability{Name: string(v.Name), Icon: string(v.Icon)}
		}
		cont = bool(q.GameData.Abilities.More)
		log.Printf("time taken on page %d: %dms. total %.1fs\n", page, time.Since(t0).Milliseconds(), time.Since(t00).Seconds())
		time.Sleep(300 * time.Millisecond)
		page++
	}

	return ret, nil
}

// TODO: GetEvents(Buffs,Friendlies) took 5.76s. Find ways to optimize oh god its slow
func (c *ClientImplementation) GetEvents(ctx context.Context,
	code string, fightID int, startTime, endTime float64, dataType EventDataType, hostilityType HostilityType) ([]*RawBuffEvent, error) {
	t0 := time.Now()
	defer func() {
		log.Printf("GetEvents(%s,%d,%.0f,%.0f,%s,%s) took %.2fs\n",
			code, fightID, startTime, endTime, dataType, hostilityType, time.Since(t0).Seconds())
	}()
	// {
	// 	reportData {
	// 		report(code: $code) {
	// 			masterData {
	// 				logVersion
	// 			}
	// 			events(
	// 				fightIDs: 27
	// 				startTime: 26641244
	// 				endTime: 27159706
	// 				dataType: Debuffs
	// 				hostilityType: Enemies
	// 			) {
	// 				nextPageTimestamp
	// 				data
	// 			}
	// 		}
	// 	}
	// }
	type query struct {
		ReportData struct {
			Report struct {
				MasterData struct {
					LogVersion graphql.Int `graphql:"logVersion"`
				} `graphql:"masterData"` // TODO: Assert
				Events struct {
					Data              stdjson.RawMessage `graphql:"data"`
					NextPageTimestamp graphql.Float      `graphql:"nextPageTimestamp"` // TODO: Paginate
				} `graphql:"events(fightIDs:[$fightID],startTime:$st,endTime:$et,dataType:$dataType,hostilityType:$hostilityType)"`
			} `graphql:"report(code:$code)"`
		} `graphql:"reportData"`
	}
	v := []*RawBuffEvent{}
	for startTime != 0 {
		q := query{}
		err := c.c.Query(ctx, &q, map[string]interface{}{
			"code":          graphql.String(code),
			"fightID":       graphql.Int(fightID),
			"st":            graphql.Float(startTime),
			"et":            graphql.Float(endTime),
			"dataType":      dataType,
			"hostilityType": hostilityType,
		})
		if err != nil {
			return nil, err
		}
		v2 := []*RawBuffEvent{}
		if err := json.Unmarshal(q.ReportData.Report.Events.Data, &v2); err != nil {
			return nil, err
		}
		startTime = float64(q.ReportData.Report.Events.NextPageTimestamp)
		v = append(v, v2...)
	}
	return v, nil
}
