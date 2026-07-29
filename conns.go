package main

import (
	"log"
	"strconv"
	"time"

	"sleeply-alive/api"
	cfg "sleeply-alive/internal/config"
	"sleeply-alive/internal/ctrl"
	. "sleeply-alive/internal/models"
	"sleeply-alive/internal/peer"
)

func startServer() {
	err := api.StartHTTPPull(":" + strconv.Itoa(cfg.Get().Port))
	if err != nil {
		log.Fatal(err)
	}
}

func clientBroadcast(cs []*peer.CPeer) {
	// 转发Msg至所有列表cpeer
	tmp := <-ctrl.NewMsgCh
	for _, p := range cs {
		select {
		case p.NewMsg <- tmp:
		default:
			log.Printf("消息阻塞:%s", p.URL)
		}
	}
}

func startClients() {
	var remotes []string = cfg.Get().Addrs
	cs := make([]*peer.CPeer, 0, len(remotes))

	for _, addr := range remotes {
		client := &peer.CPeer{URL: addr}
		cs = append(cs, client)
		// 循环启动所有CPeer,维持进程
		go func(c *peer.CPeer) {
			c.NewMsg = make(chan *Message, 1)
			defer c.Close()

			go func() {
				for {
					bowl := <-c.NewMsg
					if c.Touch {
						err := c.Send(bowl)
						if err != nil {
							log.Print(err)
						}
					}
				}
			}()

			for {
				if err := c.Connect(); err != nil {
					c.Touch = false
					log.Printf("连接异常:%v", err)
					time.Sleep(10 * time.Second)
				} else {
					log.Print("注销")
				}
			}
		}(client)
	}

	// client单push
	go func() {
		err := api.StartHTTPPush(":" + strconv.Itoa(cfg.Get().Port))
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		clientBroadcast(cs)
	}
}
