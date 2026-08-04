package ctrl

import (
	"github.com/google/uuid"
	. "sleeply-alive/internal/model"
)

var devices map[uuid.UUID]*Device = make(map[uuid.UUID]*Device)

func FindDev(id uuid.UUID) *Device {
	// 查找map对应id的*Device
	dev, exists := devices[id]
	if !exists {
		return nil
	}
	return dev
}

func MkDev(id uuid.UUID) {
	// 创建map里的*Device格
	// 有防止覆盖机制
	dev, exists := devices[id]
	if !exists { // 确认可创建记录
		dev = &Device{Status: false}
		devices[id] = dev
	}
}
