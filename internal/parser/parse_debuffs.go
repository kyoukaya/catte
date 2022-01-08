package parser

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kyoukaya/catte/internal/fflogs"
	"github.com/kyoukaya/catte/internal/xivdata"
)

// The buffs/debuffs can be found in the status exd from datamined sheets. Just add 1e6 to the ID.
var mitigDebuffs = map[int64]bool{
	1001193: true, // Reprisal
	1001195: true, // Feint
	1001203: true, // Addle
}

var mitigBuffs = map[int64]bool{
	// Tanks
	1001176: true, // Passage of Arms (Arms up)
	1001894: true, // Dark Missionary
	1001839: true, // Heart of Light
	// Healers
	1001873: true, // WHM: Temperance
	1000849: true, // AST: Collective Unconscious
	1000299: true, // SCH: Sacred Soil
	1000317: true, // SCH: Fey Illumination
	1002711: true, // SCH: Expedient (Desperate Measures)
	1002618: true, // SGE: Kerachole
	1003003: true, // SGE: Holos
	// Ranged
	1001934: true, // Troubadour
	1001951: true, // Tactician
	1001826: true, // Shield Samba
	1002707: true, // Magick Barrier
}

type DebuffEvent struct {
	ID        int64
	Start     int64
	End       int64
	RefreshTS []int64
}

type BuffEvent struct {
	ID         int64
	Start      int64
	End        int64
	PlayersHit int64
}

func (evt *DebuffEvent) String(ds *xivdata.DataSource) string {
	name := ds.Abilities[int(evt.ID)].Name
	d := time.Duration(evt.Start) * time.Millisecond
	d2 := time.Duration(evt.End) * time.Millisecond
	s := fmt.Sprintf("%02d:%02d-%02d:%02d %s",
		int(d.Minutes()), int(d.Seconds())%60,
		int(d2.Minutes()), int(d2.Seconds())%60, name)
	if len(evt.RefreshTS) != 0 {
		s += " !" + strconv.Itoa(len(evt.RefreshTS))
	}
	return s
}

func (evt *DebuffEvent) StartTs() int64 {
	return evt.Start
}
func (evt *DebuffEvent) EndTs() int64 {
	return evt.End
}
func (evt *DebuffEvent) GetID() int64 {
	return evt.ID
}
func (evt *BuffEvent) EndTs() int64 {
	return evt.End
}
func (evt *BuffEvent) GetID() int64 {
	return evt.ID
}

func (evt *BuffEvent) String(ds *xivdata.DataSource) string {
	name := ds.Abilities[int(evt.ID)].Name
	d := time.Duration(evt.Start) * time.Millisecond
	d2 := time.Duration(evt.End) * time.Millisecond
	return fmt.Sprintf("%02d:%02d-%02d:%02d %s (%d)",
		int(d.Minutes()), int(d.Seconds())%60,
		int(d2.Minutes()), int(d2.Seconds())%60, name, evt.PlayersHit)
}

func (evt *BuffEvent) StartTs() int64 {
	return evt.Start
}

func (p *Parser) ParseMitigBuff(events []*fflogs.RawBuffEvent, st int64) []*BuffEvent {
	lastBuffMap := map[int64]*BuffEvent{}
	evts := []*BuffEvent{}
	for _, v := range events {
		id := *v.AbilityGameID
		if !mitigBuffs[id] {
			continue
		}
		ts := v.Timestamp - st
		if v.Type == fflogs.Applybuff {
			last := lastBuffMap[id]
			if last == nil {
				lastBuffMap[id] = &BuffEvent{
					ID:         id,
					Start:      ts,
					PlayersHit: 1,
				}
			} else {
				last.PlayersHit++
			}
			continue
		}
		if v.Type == fflogs.Removebuff {
			last := lastBuffMap[id]
			if last == nil {
				continue
			}
			last.End = ts
			evts = append(evts, last)
			delete(lastBuffMap, id)
		}
	}
	return evts
}

func (p *Parser) ParseDebuffs(events []*fflogs.RawBuffEvent, st int64) []*DebuffEvent {
	lastDebuffMap := map[int64]*DebuffEvent{}
	evts := []*DebuffEvent{} // sorted in ascending debuff end time
	for _, v := range events {
		id := *v.AbilityGameID
		if !mitigDebuffs[id] {
			continue
		}
		ts := v.Timestamp - st
		switch v.Type {
		case fflogs.Applydebuff:
			lastDebuffMap[id] = &DebuffEvent{
				ID:    id,
				Start: ts,
			}
		case fflogs.Refreshdebuff:
			lastDebuffMap[id].RefreshTS = append(lastDebuffMap[id].RefreshTS, ts)
		case fflogs.Removedebuff:
			evt := lastDebuffMap[id]
			evt.End = ts
			evts = append(evts, evt)
		}
		if v.Type != fflogs.Removedebuff {
			continue
		}
	}
	return evts
}
