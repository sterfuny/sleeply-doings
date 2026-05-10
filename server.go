package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:	func(r *http.Request) bool {return true},
	ReadBufferSize:	1024,
	WriteBufferSize:1024,
}

//var conns = make(map[*websocket.Conn])

func handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	log.Printf("新连接: %s", conn.RemoteAddr())

	//歇息心跳
	timer := time.NewTimer(pingSpit)
	defer timer.Stop()
	
	//初次死线
	conn.SetReadDeadline(time.Now().Add(holdWait))
	//服务端设置收pong触发器
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(holdWait))
		timer.Reset(pingSpit)
		return nil
	})

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
				case <-timer.C:				
					if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
						return
					}
					timer.Reset(pingSpit)
				case <-done:
					return
			}
		}
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取失败: %v", err)
			break
		}
		conn.SetReadDeadline(time.Now().Add(holdWait))
		timer.Reset(pingSpit)

		msg, err := MessageFromJSON(msgBytes)
		if err != nil {
			log.Printf("JSON解析失败: %v", err)
			continue
		}
		log.Printf("%s", formatMessage(msg))
	}
}

func startServer(addr string) error {
	http.HandleFunc("/ws", handleConn)
	log.Printf("服务: %s", addr)
	return http.ListenAndServe(addr, nil)
}
