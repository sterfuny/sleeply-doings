package main

import "time"

const (
	writeWait	= 10 * time.Second
	holdWait	= 60 * time.Second
	pingSpit	= (holdWait*9) / 10
)
