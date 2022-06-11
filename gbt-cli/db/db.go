package db

import (
	_ "embed"
	"encoding/json"
	"github.com/thedevsaddam/gojsonq/v2"
)

//go:embed db.json
var db string

type Template struct {
	Name         string
	Artifacts    string
	Dependencies string
}

func Query(artifact string) (*Template, error) {
	j := gojsonq.New().FromString(db).From("templates").
		WhereContains("artifacts", artifact).Limit(1)
	if bb, err := json.Marshal(j.Get()); err == nil {
		t := []Template{}
		err = json.Unmarshal(bb, &t)
		return &t[0], err
	} else {
		return nil, err
	}
}
