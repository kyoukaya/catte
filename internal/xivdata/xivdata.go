package xivdata

import (
	_ "embed"

	json "github.com/json-iterator/go"
)

//go:embed ability_cache.json
var abilitiesBytes []byte

type DataSource struct {
	Abilities map[int]*Ability
}

type Ability struct {
	Name string
	Icon string
}

func NewDataSource() *DataSource {
	abilities := map[int]*Ability{}
	err := json.Unmarshal(abilitiesBytes, &abilities)
	if err != nil {
		panic(err)
	}
	return &DataSource{
		Abilities: abilities,
	}
}
