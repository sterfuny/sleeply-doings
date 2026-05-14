package main

import (
	"flag"
	"fmt"
	"log"
	"time"
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

	newMsg = make(chan struct{}, 1)
	client := &Peer{serverURL: *serverAddr}
	go func(){
		for{
			if err := client.connect(); err != nil {
				log.Printf("连接失败:%v", err)
				time.Sleep(5*time.Second)
			} else {	
				log.Print("注销")
			}
		}
	}()
	defer client.Close()

	go func(){
		if err := startHTTPPush(":9090"); err != nil {
			log.Fatal(err)
		}
	}()
	
	go func(){
		for {
			<-newMsg
			err := client.Send(msg)
			if err != nil {
				log.Print(err)
			}
		}
	}()

	log.Printf("客户端已启动，HTTP接口:9090/push，WebSocket连接: %s", *serverAddr)
	select {}
}
