package ctrl

import (
	"encoding/json"
	"fmt"
	"log"
	. "sleeply-alive/internal/models"
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
	// var parts []string
	// if v, ok := m.SendKey["AppName"].(string); ok {
	// 	parts = append(parts, fmt.Sprintf("AppName:%s", v))
	// }

	return fmt.Sprintf("%v", m.SendKey)
}

func UpdateInfo(info map[string]any) {

	var m *Message = &Message{SendKey: &info}
	select {
	case NewMsgCh <- m:
		log.Printf("HTTP推送成功:%s", FormatMessage(m))
	default:
	}
}
