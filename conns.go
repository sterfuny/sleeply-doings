package main

import (
	"time"
	"sync"
	"context"

	"github.com/gorilla/websocket"
)

const (
	writeWait	= 10 * time.Second
	holdWait	= 60 * time.Second
	pingSpit	= (holdWait*9) / 10
)

type Peer struct {
	serverURL	string
	conn		*websocket.Conn
	timer		*time.Timer
	lastMessage	*Message

	mu        sync.Mutex
	muPub     sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

var Pool map[int]Peer
