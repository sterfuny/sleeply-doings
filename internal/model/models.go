package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Message struct {
	UUID    *uuid.UUID      `json:"uuid,omitempty"`
	SendKey *map[string]any `json:"sendkey,omitempty"`
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// 仅服务端使用以自主设置,也是提供API调用的信息
type Device struct {
	Name    string   `json:"name,omitempty"`
	LastMsg *Message `json:"lastmsg,omitempty"`
	Status  bool     `json:"status"`
}

func (m *Device) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
