
package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var mode string
	port := flag.Int("port", 8080, "服务端端口")
	serverAddr := flag.String("server", "ws://localhost:8080/ws", "服务端地址(客户端模式)")
	flag.StringVar(&mode, "mode", "server", "server或client")
	flag.Parse()

	if mode == "server" {
		addr := fmt.Sprintf(":%d", *port)
		log.Fatal(startServer(addr))
	}

	// 客户端：使用参数指定的地址
	client := NewClient(*serverAddr)
	defer client.Close()

	go func() {
		if err := startHTTPPush(client, ":9090"); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("客户端已启动，HTTP接口: :9090/push，WebSocket连接: %s", *serverAddr)
	select {}
}
