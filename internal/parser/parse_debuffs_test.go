package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"

	"github.com/kyoukaya/catte/internal/fflogs"
	"github.com/kyoukaya/catte/internal/xivdata"
)

func TestParseMitigs(t *testing.T) {
	b, err := ioutil.ReadFile("EDTypeBuffsFriendlies.json")
	assert.NilError(t, err)
	teamBuffs := []*fflogs.RawBuffEvent{}
	err = json.Unmarshal(b, &teamBuffs)
	st := 26641244
	assert.NilError(t, err)
	ds := xivdata.NewDataSource()
	p := New(ds)
	evts := p.ParseMitigBuff(teamBuffs, int64(st))
	println(strings.Join(stringifyMitigs(evts, ds), "\n"))
}

func stringifyMitigs(evts []*BuffEvent, ds *xivdata.DataSource) []string {
	ret := []string{}
	for _, evt := range evts {
		name := ds.Abilities[int(evt.ID)].Name
		d := time.Duration(evt.Start) * time.Millisecond
		d2 := time.Duration(evt.End) * time.Millisecond
		ret = append(ret, fmt.Sprintf("%02d:%02d-%02d:%02d %s (%d)",
			int(d.Minutes()), int(d.Seconds())%60,
			int(d2.Minutes()), int(d2.Seconds())%60, name, evt.PlayersHit))
	}
	return ret
}

type SortByBuffEventStartAsc []*BuffEvent

func (a SortByBuffEventStartAsc) Len() int           { return len(a) }
func (a SortByBuffEventStartAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByBuffEventStartAsc) Less(i, j int) bool { return a[i].Start < a[j].Start }

func TestParseDebuffs(t *testing.T) {
	b, err := ioutil.ReadFile("EDTypeDebuffsEnemies.json")
	assert.NilError(t, err)
	enemyDebuffs := []*fflogs.RawBuffEvent{}
	err = json.Unmarshal(b, &enemyDebuffs)
	st := 26641244
	assert.NilError(t, err)
	ds := xivdata.NewDataSource()
	p := New(ds)
	evts := p.ParseDebuffs(enemyDebuffs, int64(st))
	sl := stringifyDebuffs(evts, ds)
	s := strings.Join(sl, "\n")
	println(s)
}

type SortByDebuffStart []*DebuffEvent

func (a SortByDebuffStart) Len() int           { return len(a) }
func (a SortByDebuffStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByDebuffStart) Less(i, j int) bool { return a[i].Start < a[j].Start }

func stringifyDebuffs(evts []*DebuffEvent, ds *xivdata.DataSource) []string {
	ret := []string{}
	sort.Sort(SortByDebuffStart(evts))
	for _, evt := range evts {
		name := ds.Abilities[int(evt.ID)].Name
		d := time.Duration(evt.Start) * time.Millisecond
		d2 := time.Duration(evt.End) * time.Millisecond
		ret = append(ret, fmt.Sprintf("%02d:%02d-%02d:%02d %s",
			int(d.Minutes()), int(d.Seconds())%60,
			int(d2.Minutes()), int(d2.Seconds())%60, name))
	}
	return ret
}
