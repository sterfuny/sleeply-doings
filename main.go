
package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var mode string
	port := flag.Int("port", 8080, "server port")
	serverAddr := flag.String("server", "ws://localhost:8080/ws", "server URL(client mode)")
	flag.StringVar(&mode, "mode", "server", "mode is server/client")
	// 获取用户输入参数
	flag.Parse()

	if mode == "server" {
		addr := fmt.Sprintf(":%d", *port)
		log.Fatal(startServer(addr))
	}

	client := &Peer{serverURL: *serverAddr}
	defer client.Close()

	go func() {
		if err := startHTTPPush(":9090"); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("客户端已启动，HTTP接口:9090/push，WebSocket连接: %s", *serverAddr)
	select {}
}
