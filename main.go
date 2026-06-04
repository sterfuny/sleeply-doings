package main

import (
	"log"
	"strconv"

	cfg "sleeply-alive/internal/config"
)

func main() {
	if cfg.Get().Mode == "server" {
		err := startServer(":" + strconv.Itoa(cfg.Get().Port))
		if err != nil {
			log.Println(err)
		}
	}

	if cfg.Get().Mode == "client" && cfg.Get().Addrs != nil {
		startClients(cfg.Get().Addrs)
		log.Printf("客户端已启动")
	} else {
		log.Printf("未提供Addrs,自动退出")
		return
	}
	select {}
}
