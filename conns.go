package main

import (
	"context"
	"sync"
	"time"
	"log"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 5 * time.Second
	holdWait  = 60 * time.Second
	pingSpit  = (holdWait * 9) / 10
)

type Peer struct {
	URL string
	conn      *websocket.Conn
	timer     *time.Timer
	newMsg chan *Message

	mu     sync.Mutex
	muPub  sync.RWMutex
	cancel context.CancelFunc
}

func broadcast(s []*Peer) {
	tmp := <- newMsgCh
	for _, p := range s{
		select {
		case p.newMsg <- tmp:
		default:
			log.Printf("消息阻塞:%s", p.URL)
		}
	}
}

func startClients(addrs []string) {
	Peers := make([]*Peer, 0, len(addrs))

	for _, addr := range addrs {
		client := &Peer{URL: addr}
		Peers = append(Peers, client)
		go func(c *Peer) {
			c.newMsg = make(chan *Message, 1)
			defer c.Close()
			go func() {
				for {
					err := c.Send(<-c.newMsg)
					if err != nil {
						log.Print(err)
					}
				}
			}()

			for {
				if err := c.connect(); err != nil {
					log.Printf("连接异常:%v", err)
					time.Sleep(10 * time.Second)
				} else {
					log.Print("注销")
				}
			}
		}(client)
	}

	go func() {
		if err := startHTTPPush(":9090"); err != nil {
			log.Fatal(err)
		}
	}()

	for {broadcast(Peers)}
}

