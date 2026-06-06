package manager

import (
	"fmt"
	"log"
	"strings"
	. "sleeply-alive/internal/models"
	"encoding/json"
)

var msg *Message = &Message{}
var NewMsgCh chan *Message = make(chan *Message)

func FromMessage(data []byte) (*Message, error) {
	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func FormatMessage(m *Message) string {
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

func UpdateInfo(app *AppInfo, battery *int, screen *bool) {
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

	var m *Message = &Message{App: app, Battery: battery, Screen: screen}
	select {
	case NewMsgCh <- m:
		log.Printf("HTTP推送成功:%s", FormatMessage(m))
	default:
	}
}
