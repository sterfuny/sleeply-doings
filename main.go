package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	var(
		mode string
		serverFlag string
	)
	port := flag.Int("port", 8080, "server port")
	flag.StringVar(&serverFlag , "server", "ws://localhost:8080/ws", "server URL(client mode)")
	flag.StringVar(&mode, "mode", "server", "mode is server/client")
	flag.Parse()// 获取cli输入参数

	serverAddrs := strings.Split(serverFlag, ",")

	if mode == "server" {
		addr := fmt.Sprintf(":%d", *port)
		log.Fatal(startServer(addr))
	}

	// mode = "debug"
	if mode == "debug" {
		serverAddrs = append(serverAddrs, 
			"ws://localhost:8181/ws",
			"ws://localhost:9191/ws",
		)
	}

	startClients(serverAddrs)

	log.Printf("客户端已启动")
	select {}
}
