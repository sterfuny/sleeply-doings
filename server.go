package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("升级失败: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("新连接: %s", conn.RemoteAddr())

	//歇息心跳
	idleTimeout := 30 * time.Second
	timer := time.NewTimer(idleTimeout)
	defer timer.Stop()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		timer.Reset(idleTimeout)
		return nil
	})

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <-timer.C:
				if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)); err != nil {
					return
				}
				timer.Reset(idleTimeout)
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
    conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // 每条消息都重置
    timer.Reset(30 * time.Second)

    msg, err := MessageFromJSON(msgBytes)
    if err != nil {
        log.Printf("JSON解析失败: %v", err)
        continue
    }
    log.Printf("收到: %s", formatMessage(msg))
	}
}

func startServer(addr string) error {
	http.HandleFunc("/ws", handleConn)
	log.Printf("服务启动在 %s", addr)
	return http.ListenAndServe(addr, nil)
}