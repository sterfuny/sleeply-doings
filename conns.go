package main

import (
	"log"
	"strconv"
	"time"

	cfg "sleeply-alive/internal/config"
	. "sleeply-alive/internal/models"
	"sleeply-alive/internal/peer"
)

func startServer(addr string) error {
	return peer.RegisterHandle(addr)
}

func clientBroadcast(s []*peer.CPeer) {
	tmp := <-NewMsgCh
	for _, p := range s {
		select {
		case p.NewMsg <- tmp:
		default:
			log.Printf("消息阻塞:%s", p.URL)
		}
	}
}

func startClients(addrs []string) {
	Peers := make([]*peer.CPeer, 0, len(addrs))

	for _, addr := range addrs {
		client := &peer.CPeer{URL: addr}
		Peers = append(Peers, client)
		go func(c *peer.CPeer) {
			c.NewMsg = make(chan *Message, 1)
			defer c.Close()
			go func() {
				for {
					err := c.Send(<-c.NewMsg)
					if err != nil {
						log.Print(err)
					}
				}
			}()

			for {
				if err := c.Connect(); err != nil {
					log.Printf("连接异常:%v", err)
					time.Sleep(10 * time.Second)
				} else {
					log.Print("注销")
				}
			}
		}(client)
	}

	go func() {
		if err := startHTTPPush(":"+strconv.Itoa(cfg.Get().Port)); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		clientBroadcast(Peers)
	}
}
