package parser

import (
	"github.com/kyoukaya/catte/internal/fflogs"
	"github.com/kyoukaya/catte/internal/xivdata"
)

type EventInfo struct {
	ID          int64
	Ts          int64
	OrigTs      int64
	TotalDamage int64
	Instances   int
	SubInfo     []*EventInfo // for fast hits, consecutive
}

type Parser struct {
	attackIDs map[int64]struct{}
}

func New(ds *xivdata.DataSource) *Parser {
	m := map[int64]struct{}{}
	for id, a := range ds.Abilities {
		if a.Name != "attack" {
			continue
		}
		m[int64(id)] = struct{}{}
	}
	return &Parser{
		attackIDs: m,
	}
}

func (p *Parser) ParseDamageTaken(events []*fflogs.RawBuffEvent, st int64) []*EventInfo {
	type idts struct {
		id int64
		ts int64
	}
	seen := map[idts]*EventInfo{}
	tl := []idts{}
	for _, event := range events {
		if event.Type != fflogs.Calculateddamage {
			continue
		}
		if v := seen[idts{id: *event.AbilityGameID, ts: event.Timestamp}]; v != nil {
			v.Instances++
			v.TotalDamage += *event.Amount
			continue
		}
		seen[idts{id: *event.AbilityGameID, ts: event.Timestamp}] = &EventInfo{
			ID:          *event.AbilityGameID,
			Ts:          event.Timestamp - st,
			TotalDamage: *event.Amount,
			Instances:   1,
		}
		tl = append(tl, idts{id: *event.AbilityGameID, ts: event.Timestamp})
	}
	stack := []*EventInfo{}
	infoTL := []*EventInfo{}
	for i, hash := range tl {
		v := seen[hash]
		// Close hit compression
		var isCloseHit bool
		var shouldCompressAttack bool
		if i > 0 {
			prev := seen[tl[i-1]]
			if v.ID == prev.ID {
				isCloseHit = v.Ts-prev.Ts < 1e3
				_, isAttack := p.attackIDs[v.ID]
				shouldCompressAttack = isAttack
			}
		}
		// Attack compression
		if shouldCompressAttack || isCloseHit {
			if len(stack) == 0 {
				stack = append(stack, seen[tl[i-1]])
				infoTL = infoTL[:len(infoTL)-1]
			}
			stack = append(stack, v)
			continue
		}
		if n := len(stack); n > 0 {
			v2 := stack[0]
			infoTL = append(infoTL, v2)
			v2.SubInfo = stack[0:]
			stack = nil
		}
		infoTL = append(infoTL, v)
	}
	// pop off any close hit compression, attack compression
	if len(stack) > 0 {
		v2 := stack[0]
		infoTL = append(infoTL, v2)
		v2.SubInfo = stack[0:]
		stack = nil
	}
	return infoTL
}
