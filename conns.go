package main

import (
	"time"
	"sync"
	"context"

	"github.com/gorilla/websocket"
)

const (
	writeWait	= 5 * time.Second
	holdWait	= 60 * time.Second
	pingSpit	= (holdWait*9) / 10
)

type Peer struct {
	serverURL	string
	conn		*websocket.Conn
	timer		*time.Timer
	lastMsg		*Message

	mu        sync.Mutex
	muPub     sync.RWMutex
	cancel    context.CancelFunc
}

var Pool map[int]*Peer
var Poolindex int = 0
/*
func addPeer(p *Peer) int{
	Poolindex++
	Pool[Poolindex] = p
	return Poolindex
}
*/
