package manager

import (
	. "sleeply-alive/internal/models"
)

var devices map[string]*Device = make(map[string]*Device)

func FindDev(id string) *Device{
	dev, exists := devices[id]
	if !exists { // 新设备,创建记录
		dev = &Device{}
		devices[id] = dev
	}
	return devices[id]
}
