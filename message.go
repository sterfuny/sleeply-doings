package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 当前应用,电量,屏幕 握手的双方都可用
// 如果是单个请求写url参数是可以的,但需要长连接还是用请求体好
type Message struct {
	App     *AppInfo `json:"app,omitempty"`
	Battery *int     `json:"battery,omitempty"`
	Screen  *bool    `json:"screen,omitempty"`
}

type AppInfo struct {
	Name string `json:"name,omitempty"`
	Pkg  string `json:"pkg,omitempty"`
}

// 仅服务端使用以自主设置,也是提供API调用的信息
type Device struct {
	ID		string  	`json:"id"`
	Name	string		`json:"name,omitempty"`
	Lastmsg	*Message	`json:"lastmsg,omitempty"`
	Status	bool		`json:"status"`
}

var msg *Message = &Message{}
var newMsgCh chan *Message = make(chan *Message)

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Device) ToJSON() ([]byte, error) {
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

func updateInfo(app *AppInfo, battery *int, screen *bool) {
	if app != nil {
		if msg.App == nil {
			msg.App = app
		} else {
			if app.Name != "" {
				msg.App.Name = app.Name
			}
			if app.Pkg != "" {
				msg.App.Pkg = app.Pkg
			}
		}
	}
	if battery != nil {
		msg.Battery = battery
	}
	if screen != nil {
		msg.Screen = screen
	}

	select {
	case newMsgCh <- &Message{App: app, Battery: battery, Screen: screen}:
	default:
	}
}
