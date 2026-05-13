package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

//如果是单个请求写url参数是可以的,但需要长连接还是用请求体好
type Message struct {
	App     *AppInfo `json:"app,omitempty"`
	Battery *int     `json:"battery,omitempty"`
	Screen  *bool    `json:"screen,omitempty"`
}

type AppInfo struct {
	Name string `json:"name,omitempty"`
	Pkg  string `json:"pkg,omitempty"`
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func MessageFromJSON(data []byte) (*Message, error) {
	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func formatMessage(m *Message) string {
	var parts []string
	if m.App != nil {
		parts = append(parts, fmt.Sprintf("App:{Name:%s Pkg:%s}", m.App.Name, m.App.Pkg))
	}
	if m.Battery != nil {
		parts = append(parts, fmt.Sprintf("Battery:%d", *m.Battery))
	}
	if m.Screen != nil {
		parts = append(parts, fmt.Sprintf("Screen:%t", *m.Screen))
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, " "))
}
