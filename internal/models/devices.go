package models

import (
	"encoding/json"
)

// 仅服务端使用以自主设置,也是提供API调用的信息
type Device struct {
	// ID      string   `json:"id"`
	Name    string   `json:"name,omitempty"`
	Lastmsg *Message `json:"lastmsg,omitempty"`
	Status  bool     `json:"status"`
}

func (m *Device) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
