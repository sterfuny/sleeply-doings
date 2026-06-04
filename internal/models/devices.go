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

var devices map[string]*Device = make(map[string]*Device)

func Find(id string) *Device{
	dev, exists := devices[id]
	if !exists { // 新设备,创建记录
		dev = &Device{}
		devices[id] = dev
	}
	return devices[id]
}
