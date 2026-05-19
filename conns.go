package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 5 * time.Second
	holdWait  = 60 * time.Second
	pingSpit  = (holdWait * 9) / 10
)

type sPeer struct {
	divice *Device 
	conn   *websocket.Conn
	timer  *time.Timer

	ctx    context.Context
	cancel context.CancelFunc
}

type cPeer struct {
	URL    string
	conn   *websocket.Conn
	newMsg chan *Message

	mu    sync.Mutex
	muPub sync.RWMutex
}

var devices map[string]*Device

func startServer(addr string) error {
	http.HandleFunc("/ws", handleConn)
	log.Printf("服务启动:%s", addr)
	return http.ListenAndServe(addr, nil)
}

func clientBroadcast(s []*cPeer) {
	tmp := <-newMsgCh
	for _, p := range s {
		select {
		case p.newMsg <- tmp:
		default:
			log.Printf("消息阻塞:%s", p.URL)
		}
	}
}

func startClients(addrs []string) {
	Peers := make([]*cPeer, 0, len(addrs))

	for _, addr := range addrs {
		client := &cPeer{URL: addr}
		Peers = append(Peers, client)
		go func(c *cPeer) {
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
		if err := startHTTPPush(strconv.Itoa(Get().Port)); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		clientBroadcast(Peers)
	}
}
