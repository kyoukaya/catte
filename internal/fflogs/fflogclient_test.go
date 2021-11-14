package fflogs

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"gotest.tools/assert"
)

type info struct {
	Name        string
	Ts          int64
	OrigTs      int64
	TotalDamage int64
	Instances   int
	subInfo     []*info // for fast hits, consecutive
}

func TestNewClient(t *testing.T) {
	got, err := NewClient("94df231e-cd7b-4658-b051-90cfb379db49",
		"To6WOQzrGxSamFfh6oUKsGxupVpYH1kI0qBAnlNT",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiI5NGRmMjMxZS1jZDdiLTQ2NTgtYjA1MS05MGNmYjM3OWRiNDkiLCJqdGkiOiIwODI4ZTdhNWNhZGJiZmY0ZDlhMDE1MTMxMDBiYTBmMDc1MTNiNjRlN2Q3NjA0NzM3YTEyNDNhYWQ2Zjg0N2JhZWI2Y2UzMDk3MGY0ZjJkZSIsImlhdCI6MTYzNjg2ODgyNywibmJmIjoxNjM2ODY4ODI3LCJleHAiOjE2NDcyMzY4MjcsInN1YiI6IiIsInNjb3BlcyI6WyJ2aWV3LXVzZXItcHJvZmlsZSIsInZpZXctcHJpdmF0ZS1yZXBvcnRzIl19.vmeK4w3rHbx1bjm-4h1C89VoebdUpGuIcR9djmFGfWoyoUanOJqhw4FupY2DUykhPlcaX9QO9zZOal3ZeZunM9m-9-5rBlUffCa9kpJadA1ysADrW-MnQ0UDQMbudUxUnuhKXS7FHZNBbgYRPQNXQfHbK0msLqx7F6GhlQmyiXvCfJ1Ol8lS594tZeOuQZzyjf0yHGw6CpSAAHwm0oW51H69EmXkzwNA0-HgdNlTGCcPdGL4Ln20cXK2I8RbYhUbDjlwe_GeV8xeSh6EI6aN-1xr05RP1DIpW5MWaWFliQQPdO6JJAdv-7lGU9CdKMhwW9JhEc3aXMnD0GmLQYT_BupFoItJLxy5ZyX5EKrdWA0G-qOGQ1HfqhkZLSGO0leuIrbnqFqeJfiRNUEi2d3_Ot-9tX59ikBX8CNTgEYGfRsm2bPSMgyZe3ygdOO17wTWNFuJVBdgdIH12wHL_Lt1OMxCQogA9GgvFfOuA6cqcETJ6HMbOnlbbacCAQia8yHAT6eyoAckYqRjZOAwA4AVoctwX-SXaBRu-yICIdEr1DigYFRVoen6Xc8xu0B7YiOYnXRhEGKCJPYm8I3v5llidfPTlMlSWTAsdzYO38CFW2YvYnSue4qTjZK1gGY7CEPJEs-vY0bWgASJxNXuDHOn8BRLW405ZPQptkVhuDwGqGQ")
	assert.NilError(t, err)
	ctx := context.Background()
	fightID := 27
	code := "6bPT1tMkVLjrZxwN"
	st, et, err := got.GetTimesFromFightAndID(ctx, code, fightID)
	assert.NilError(t, err)
	v, err := got.GetEvents(ctx, code, fightID, st, et, EDTypeDebuffs, Enemies)
	assert.NilError(t, err)
	unmarshalSave(t, "EDTypeDebuffsEnemies", v)
	v, err = got.GetEvents(ctx, code, fightID, st, et, EDTypeBuffs, Friendlies)
	assert.NilError(t, err)
	unmarshalSave(t, "EDTypeBuffsFriendlies", v)
	v, err = got.GetEvents(ctx, code, fightID, st, et, EDTypeDamageTaken, Friendlies)
	assert.NilError(t, err)
	unmarshalSave(t, "EDTypeDamageTakenFriendlies", v)
	v, err = got.GetEvents(ctx, code, fightID, st, et, EDTypeHealing, Friendlies)
	assert.NilError(t, err)
	unmarshalSave(t, "EDTypeHealingFriendlies", v)
	v, err = got.GetEvents(ctx, code, fightID, st, et, EDTypeCasts, Friendlies)
	unmarshalSave(t, "EDTypeCastsFriendlies", v)
	assert.NilError(t, err)
}

func unmarshalSave(t *testing.T, filename string, v interface{}) {
	b, err := json.Marshal(v)
	assert.NilError(t, err)
	err = ioutil.WriteFile(filename, b, 0600)
	assert.NilError(t, err)
}

func TestPaginate(t *testing.T) {
	cli, err := NewClient("94df231e-cd7b-4658-b051-90cfb379db49",
		"To6WOQzrGxSamFfh6oUKsGxupVpYH1kI0qBAnlNT",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiI5NGRmMjMxZS1jZDdiLTQ2NTgtYjA1MS05MGNmYjM3OWRiNDkiLCJqdGkiOiIwODI4ZTdhNWNhZGJiZmY0ZDlhMDE1MTMxMDBiYTBmMDc1MTNiNjRlN2Q3NjA0NzM3YTEyNDNhYWQ2Zjg0N2JhZWI2Y2UzMDk3MGY0ZjJkZSIsImlhdCI6MTYzNjg2ODgyNywibmJmIjoxNjM2ODY4ODI3LCJleHAiOjE2NDcyMzY4MjcsInN1YiI6IiIsInNjb3BlcyI6WyJ2aWV3LXVzZXItcHJvZmlsZSIsInZpZXctcHJpdmF0ZS1yZXBvcnRzIl19.vmeK4w3rHbx1bjm-4h1C89VoebdUpGuIcR9djmFGfWoyoUanOJqhw4FupY2DUykhPlcaX9QO9zZOal3ZeZunM9m-9-5rBlUffCa9kpJadA1ysADrW-MnQ0UDQMbudUxUnuhKXS7FHZNBbgYRPQNXQfHbK0msLqx7F6GhlQmyiXvCfJ1Ol8lS594tZeOuQZzyjf0yHGw6CpSAAHwm0oW51H69EmXkzwNA0-HgdNlTGCcPdGL4Ln20cXK2I8RbYhUbDjlwe_GeV8xeSh6EI6aN-1xr05RP1DIpW5MWaWFliQQPdO6JJAdv-7lGU9CdKMhwW9JhEc3aXMnD0GmLQYT_BupFoItJLxy5ZyX5EKrdWA0G-qOGQ1HfqhkZLSGO0leuIrbnqFqeJfiRNUEi2d3_Ot-9tX59ikBX8CNTgEYGfRsm2bPSMgyZe3ygdOO17wTWNFuJVBdgdIH12wHL_Lt1OMxCQogA9GgvFfOuA6cqcETJ6HMbOnlbbacCAQia8yHAT6eyoAckYqRjZOAwA4AVoctwX-SXaBRu-yICIdEr1DigYFRVoen6Xc8xu0B7YiOYnXRhEGKCJPYm8I3v5llidfPTlMlSWTAsdzYO38CFW2YvYnSue4qTjZK1gGY7CEPJEs-vY0bWgASJxNXuDHOn8BRLW405ZPQptkVhuDwGqGQ")
	assert.NilError(t, err)
	ctx := context.Background()
	v, err := cli.GetAllAbilities(ctx)
	assert.NilError(t, err)
	b, err := json.Marshal(v)
	assert.NilError(t, err)
	ioutil.WriteFile("ability_cache", b, 0600)
}
