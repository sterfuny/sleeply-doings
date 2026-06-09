package main

import (
	"log"
	"strconv"
	"time"

	cfg "sleeply-alive/internal/config"
	. "sleeply-alive/internal/models"
	"sleeply-alive/internal/peer"
	"sleeply-alive/internal/ctrl"
	"sleeply-alive/api"
)

func startServer(addr string) error {
	go func() {
		if err := api.StartHTTPPull(addr); err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

func clientBroadcast(s []*peer.CPeer) {
	tmp := <-ctrl.NewMsgCh
	for _, p := range s {
		select {
		case p.NewMsg <- tmp:
		default:
			log.Printf("消息阻塞:%s", p.URL)
		}
	}
}

func startClients(remotes []string) {
	Peers := make([]*peer.CPeer, 0, len(remotes))

	for _, addr := range remotes {
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
		if err := api.StartHTTPPush(":"+strconv.Itoa(cfg.Get().Port)); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		clientBroadcast(Peers)
	}
}
