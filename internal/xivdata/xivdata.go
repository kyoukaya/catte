package xivdata

import (
	_ "embed"
	"encoding/csv"
	"net/http"
	"strconv"
)

type DataSource struct {
	Abilities map[int]*Ability
}

type Ability struct {
	Name string
}

func NewDataSource() (*DataSource, error) {
	abilities := map[int]*Ability{}
	resp, err := http.Get("https://raw.githubusercontent.com/xivapi/ffxiv-datamining/master/csv/Action.csv")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r := csv.NewReader(resp.Body)
	for i := 0; i < 3; i++ {
		r.Read()
	}
	for l, err := r.Read(); err == nil; l, err = r.Read() {
		i, err := strconv.Atoi(l[0])
		if err != nil {
			return nil, err
		}
		abilities[i] = &Ability{Name: l[1]}
	}
	return &DataSource{
		Abilities: abilities,
	}, nil
}
