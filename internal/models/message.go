package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Message struct {
	// App     *AppInfo `json:"app,omitempty"`
	// Battery *int     `json:"battery,omitempty"`
	// Screen  *bool    `json:"screen,omitempty"`
	UUID	*uuid.UUID		`json:"uuid,omitempty"`
	SendKey	*map[string]any	`json:"sendkey,omitempty"`
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
