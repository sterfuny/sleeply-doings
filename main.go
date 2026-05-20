package main

import (
	"fmt"
	"log"
	"strconv"
)

func Get() *Config {
	return config
}

func main() {
	if err := configInit(); err != nil {
		fmt.Println("Failed to load configuration")
		panic(err)
	}

	if Get().Mode == "server" {
		err := startServer(":" + strconv.Itoa(Get().Port))
		if err != nil {
			log.Println(err)
		}
	}

	if Get().Mode == "client" && Get().Addrs != nil {
		startClients(Get().Addrs)
		log.Printf("客户端已启动")
	} else {
		log.Printf("未提供Addrs,自动退出")
		return
	}
	select {}
}
