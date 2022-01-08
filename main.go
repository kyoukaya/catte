package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path"
	"runtime/debug"
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

type DiscordSession interface {
	AddHandler(handler interface{}) func()
	ApplicationCommandCreate(appID string, guildID string, cmd *discordgo.ApplicationCommand) (ccmd *discordgo.ApplicationCommand, err error)
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) (err error)
	FollowupMessageCreate(appID string, interaction *discordgo.Interaction, wait bool, data *discordgo.WebhookParams) (*discordgo.Message, error)
}

type App struct {
	userID          string
	xivds           *xivdata.DataSource
	fflogs          fflogs.Client
	parser          *parser.Parser
	commandHandlers map[string]Handler
}
type Handler func(ctx context.Context, s DiscordSession, i *discordgo.InteractionCreate) error
type Command struct {
	*discordgo.ApplicationCommand
	f Handler
}

const envPrefix = "CATTE_"

func main() {
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
	ds, err := xivdata.NewDataSource()
	check(err)
	app := &App{
		xivds:  ds,
		fflogs: fflogsClient,
		parser: parser.New(ds),
	}

	dg, err := discordgo.New("Bot " + token)
	check(err)
	err = dg.Open()
	if err != nil {
		panic(err)
	}
	app.RegisterSlashCommands(dg.State.User.ID, dg)
	defer dg.Close()

	// Wait here until CTRL-C or other term signal is received.
	log.Printf("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (a *App) RegisterSlashCommands(userID string, dg DiscordSession) {
	a.userID = userID
	dg.AddHandler(a.slashWrap)
	commands := []*Command{
		{
			ApplicationCommand: &discordgo.ApplicationCommand{
				Name:        "mitig",
				Description: "Get mitigation information from an FFLogs fight",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "fflogs-url",
						Description: "FFLogs URL",
						Required:    true,
					},
				},
			},
			f: a.MitigHandler,
		},
		{
			ApplicationCommand: &discordgo.ApplicationCommand{
				Name:        "dmgin",
				Description: "Get incoming damage information from an FFLogs fight",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "fflogs-url",
						Description: "FFLogs URL",
						Required:    true,
					},
				},
			},
			f: a.DamageInHandler,
		},
	}
	for _, v := range commands {
		_, err := dg.ApplicationCommandCreate(userID, "", v.ApplicationCommand)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}
	a.commandHandlers = map[string]Handler{}
	for _, command := range commands {
		log.Print("registered slash command handler ", command.Name)
		a.commandHandlers[command.Name] = command.f
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (a *App) slashWrap(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := i.ApplicationCommandData().Name
	f := a.commandHandlers[cmd]
	if f == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] panic:\n" + string(debug.Stack()))
		}
	}()
	t0 := time.Now()
	defer func() { log.Printf("command %q took %.2fs", cmd, time.Since(t0).Seconds()) }()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := f(ctx, s, i)
	if err != nil {
		msg := fmt.Sprintf("command %q failed: %v", cmd, err)
		log.Printf("[ERROR] %v", msg)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: msg},
		})
	}
}

func (a *App) DamageInHandler(ctx context.Context, s DiscordSession, i *discordgo.InteractionCreate) error {
	t0 := time.Now()
	input := i.ApplicationCommandData().Options[0].StringValue()
	events, err := a.GetDamageCommand(ctx, input)
	if err != nil {
		return err
	}
	dmginOutput := strings.Join(stringifyAttacks(a.xivds, events), "\n")
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("/stretch'd for %.2fs", time.Since(t0).Seconds()),
			Files:   []*discordgo.File{{ContentType: "text/plain", Name: "damagein.txt", Reader: strings.NewReader(dmginOutput)}},
		},
	})
}

func (a *App) MitigHandler(ctx context.Context, s DiscordSession, i *discordgo.InteractionCreate) error {
	t0 := time.Now()
	input := i.ApplicationCommandData().Options[0].StringValue()
	// Should auto switch over to followup message ideally for all messages
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return err
	}
	evts, err := a.GetMitigUsage(ctx, input)
	if err != nil {
		return err
	}
	mitigOutput := strings.Join(stringifyBuffsDebuffs(a.xivds, evts), "\n")
	m, err := s.FollowupMessageCreate(a.userID, i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("/stretch'd for %.2fs", time.Since(t0).Seconds()),
		Files:   []*discordgo.File{{ContentType: "text/plain", Name: "damagein.txt", Reader: strings.NewReader(mitigOutput)}},
	})
	log.Print(m.Content)
	return err
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
		name := xivds.Actions[int(v.ID)].Name
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
	log.Printf("evts: %#v", evts)
	for _, evt := range evts {
		ret = append(ret, evt.String(ds))
	}
	return ret
}
