package models

import (
	"encoding/json"
)

type AppInfo struct {
	Name string `json:"name,omitempty"`
	Pkg  string `json:"pkg,omitempty"`
}

type Message struct {
	App     *AppInfo `json:"app,omitempty"`
	Battery *int     `json:"battery,omitempty"`
	Screen  *bool    `json:"screen,omitempty"`
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
