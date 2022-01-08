package xivdata

import (
	_ "embed"
	"encoding/csv"
	"net/http"
	"strconv"
)

type DataSource struct {
	Actions  map[int]*Action
	Statuses map[int]*Status
}

type Action struct {
	Name string
}

type Status struct {
	Name string
}

func NewDataSource() (*DataSource, error) {
	actions := map[int]*Action{}
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
		actions[i] = &Action{Name: l[1]}
	}
	statuses := map[int]*Status{}
	resp, err = http.Get("https://raw.githubusercontent.com/xivapi/ffxiv-datamining/master/csv/Status.csv")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r = csv.NewReader(resp.Body)
	for i := 0; i < 3; i++ {
		r.Read()
	}
	for l, err := r.Read(); err == nil; l, err = r.Read() {
		i, err := strconv.Atoi(l[0])
		if err != nil {
			return nil, err
		}
		statuses[i] = &Status{Name: l[1]}
	}
	return &DataSource{
		Actions:  actions,
		Statuses: statuses,
	}, nil
}
