package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/kyoukaya/catte/internal/fflogs"
	"github.com/kyoukaya/catte/internal/parser"
	"github.com/kyoukaya/catte/internal/xivdata"
)

const prefix = ","

type App struct {
	xivds  *xivdata.DataSource
	fflogs fflogs.Client
	parser *parser.Parser
}

func main() {
	const envPrefix = "CATTE_"
	var (
		token          = os.Getenv(envPrefix + "DISCORDTOKEN")
		fflogsClientID = os.Getenv(envPrefix + "FFLOGSID")
		fflogsToken    = os.Getenv(envPrefix + "FFLOGSTOKEN")
	)
	f, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	check(err)
	defer f.Close()
	log.SetOutput(f)
	fflogsClient, err := fflogs.NewClient(fflogsClientID, fflogsToken)
	check(err)
	ds := xivdata.NewDataSource()
	app := &App{
		xivds:  ds,
		fflogs: fflogsClient,
		parser: parser.New(ds),
	}

	dg, err := discordgo.New("Bot " + token)
	check(err)
	dg.AddHandler(app.MessageHandler)
	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Printf("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (a *App) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}
	err := a.router(s, m)
	if err != nil {
		log.Printf("command failed: %v", err)
		_, _ = s.ChannelMessageSend(m.ChannelID, "command failed: "+err.Error())
	}
}

func (a *App) router(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sl := strings.SplitN(strings.TrimPrefix(m.Content, prefix), " ", 2)
	cmd := string(sl[0])
	t0 := time.Now()
	defer func() { log.Printf("command %q took %.2fs", cmd, time.Since(t0).Seconds()) }()
	switch cmd {
	case "dmgin":
		if len(sl) > 1 {
			t0 := time.Now()
			events, err := a.GetDamageCommand(ctx, string(sl[1]))
			if err != nil {
				return err
			}
			msg := strings.Join(stringifyAttacks(a.xivds, events), "\n")
			_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Reference: &discordgo.MessageReference{
					MessageID: m.ID,
					ChannelID: m.ChannelID,
					GuildID:   m.GuildID,
				},
				Content: fmt.Sprintf("/stretch'd for %.2fs", time.Since(t0).Seconds()),
				Files: []*discordgo.File{{
					Name:   "timeline.txt",
					Reader: strings.NewReader(msg),
				}},
			})
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("usage %sdmgin {fflogs URL}", prefix)
		}
	case "mitig":
		if len(sl) > 1 {
			t0 := time.Now()
			go func() { _ = s.ChannelTyping(m.ChannelID) }()
			evts, err := a.GetMitigUsage(ctx, string(sl[1]))
			if err != nil {
				return err
			}
			msg := strings.Join(stringifyBuffsDebuffs(a.xivds, evts), "\n")
			_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Reference: &discordgo.MessageReference{
					MessageID: m.ID,
					ChannelID: m.ChannelID,
					GuildID:   m.GuildID,
				},
				Content: fmt.Sprintf("/stretch'd for %.2fs", time.Since(t0).Seconds()),
				Files: []*discordgo.File{{
					Name:   "mitigs.txt",
					Reader: strings.NewReader(msg),
				}},
			})
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("usage %smitig {fflogs URL}", prefix)
		}
	default:
		return fmt.Errorf("not sure what you meant, try asking another c@")
	}
	return nil
}

func (a *App) GetMitigUsage(ctx context.Context, cmd string) ([]buffdebuffInterface, error) {
	cmd = strings.TrimSpace(cmd)
	code, fightID, err := getFightFromURL(cmd)
	if err != nil {
		// TODO: fallback to other inputs
		return nil, err
	}
	st, et, err := a.fflogs.GetTimesFromFightAndID(ctx, code, fightID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fight timings: %w", err)
	}
	// TODO: concurrency
	ret := []buffdebuffInterface{}
	debuffs, err := a.fflogs.GetEvents(ctx, code, fightID, st, et, fflogs.EDTypeDebuffs, fflogs.Enemies)
	if err != nil {
		return nil, fmt.Errorf("failed to get debuff from fflogs: %w", err)
	}
	for _, v := range a.parser.ParseDebuffs(debuffs, int64(st)) {
		ret = append(ret, v)
	}
	buffs, err := a.fflogs.GetEvents(ctx, code, fightID, st, et, fflogs.EDTypeBuffs, fflogs.Friendlies)
	if err != nil {
		return nil, fmt.Errorf("failed to get buffs from fflogs: %w", err)
	}
	for _, v := range a.parser.ParseMitigBuff(buffs, int64(st)) {
		ret = append(ret, v)
	}
	return ret, err
}

func (a *App) GetDamageCommand(ctx context.Context, cmd string) ([]*parser.EventInfo, error) {
	cmd = strings.TrimSpace(cmd)
	code, fightID, err := getFightFromURL(cmd)
	if err != nil {
		// TODO: fallback to other inputs
		return nil, err
	}
	st, et, err := a.fflogs.GetTimesFromFightAndID(ctx, code, fightID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fight timings: %w", err)
	}
	// TODO: This should be concurrent
	dmgEvents, err := a.fflogs.GetEvents(ctx, code, fightID, st, et, fflogs.EDTypeDamageTaken, fflogs.Friendlies)
	if err != nil {
		return nil, fmt.Errorf("failed to get damage events from fflogs: %w", err)
	}
	// TODO:
	// buffs, err := a.fflogs.GetEvents(ctx, code, fightID, st, et, fflogs.EDTypeBuffs, fflogs.Friendlies)
	// if err != nil {
	// 	return fmt.Errorf("failed to get damage events from fflogs: %w", err)
	// }
	// debuffs, err := a.fflogs.GetEvents(ctx, code, fightID, st, et, fflogs.EDTypeDebuffs, fflogs.Enemies)
	// if err != nil {
	// 	return fmt.Errorf("failed to get damage events from fflogs: %w", err)
	// }
	return parser.New(a.xivds).ParseDamageTaken(dmgEvents, int64(st)), nil
}

func getFightFromURL(cmd string) (code string, fightID int, err error) {
	u, err := url.Parse(cmd)
	if err != nil {
		return "", 0, fmt.Errorf("invalid url: %w", err)
	}
	switch u.Hostname() {
	case "www.fflogs.com":
	case "fflogs.com":
	default:
		return "", 0, fmt.Errorf("invalid url: unknown hostname %q", u.Hostname())
	}
	dir, tail := path.Split(u.Path)
	if dir != "/reports/" {
		return "", 0, fmt.Errorf("invalid url: invalid path %q", u.Path)
	}
	v, err := url.ParseQuery(u.Fragment)
	if err != nil {
		return "", 0, fmt.Errorf("invalid url: invalid fight fragment")
	}
	fightNumStr := v.Get("fight")
	if fightNumStr == "" {
		return "", 0, fmt.Errorf("invalid url: missing fight fragment")
	}
	fightID, err = strconv.Atoi(fightNumStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid url: invalid fight fragment")
	}
	return tail, fightID, nil
}

func stringifyAttacks(xivds *xivdata.DataSource, infoTL []*parser.EventInfo) []string {
	res := []string{}
	for _, v := range infoTL {
		name := xivds.Abilities[int(v.ID)].Name
		if name == "attack" {
			name = "AA"
		}
		if len(v.SubInfo) == 0 {
			d := (time.Duration(v.Ts)) * time.Millisecond
			res = append(res, fmt.Sprintf("%02d:%02d %s: %.0f (%d)",
				int(d.Minutes()), int(d.Seconds())%60,
				name, float64(v.TotalDamage)/float64(v.Instances), v.Instances))
			continue
		}
		stack := v.SubInfo
		n := len(stack)
		v2 := stack[0]
		d := (time.Duration(v2.Ts)) * time.Millisecond
		vmax := (time.Duration(stack[n-1].Ts)) * time.Millisecond
		tdmg, instance := sumInfoDamages(stack)
		res = append(res, fmt.Sprintf("%02d:%02d %s×%d in %ds: %.0f (%d)",
			int(d.Minutes()), int(d.Seconds())%60, name, n,
			(vmax.Milliseconds()-v2.Ts)/1e3, float64(tdmg)/float64(instance), instance/len(stack)))
	}
	return res
}

func sumInfoDamages(is []*parser.EventInfo) (totalDamage int64, instances int) {
	for _, i := range is {
		totalDamage += i.TotalDamage
		instances += i.Instances
	}
	return
}

type SortByDebuffStart []buffdebuffInterface

func (a SortByDebuffStart) Len() int           { return len(a) }
func (a SortByDebuffStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByDebuffStart) Less(i, j int) bool { return a[i].StartTs() < a[j].StartTs() }

type buffdebuffInterface interface {
	String(ds *xivdata.DataSource) string
	StartTs() int64
	EndTs() int64
	GetID() int64
}

func stringifyBuffsDebuffs(ds *xivdata.DataSource, evts []buffdebuffInterface) []string {
	ret := []string{}
	sort.Sort(SortByDebuffStart(evts))
	for _, evt := range evts {
		ret = append(ret, evt.String(ds))
	}
	return ret
}
