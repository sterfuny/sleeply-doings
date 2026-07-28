package ctrl

import (
	. "sleeply-alive/internal/models"
)

var devices map[string]*Device = make(map[string]*Device)

func FindDev(id string) *Device {
	// 查找map对应id的*Device
	dev, exists := devices[id]
	if !exists {
		return nil
	}
	return dev
}

func MkDev(id string) {
	// 创建map里的*Device格
	// 有防止覆盖机制
	dev, exists := devices[id]
	if !exists { // 确认可创建记录
		dev = &Device{}
		devices[id] = dev
	}
}
