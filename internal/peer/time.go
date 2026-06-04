package peer

import (
	"time"
)

const (
	writeWait = 5 * time.Second
	holdWait  = 60 * time.Second
	pingSpit  = (holdWait * 9) / 10
)
