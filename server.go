package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	ws := &sPeer{conn: conn}
	go ws.handleWSClient()
}

func (s *sPeer) handleWSClient() {
	conn := s.conn
	defer conn.Close()

	log.Printf("新连接:%s", conn.RemoteAddr())

	// 放入计时
	timer := time.NewTimer(pingSpit)
	s.timer = timer
	defer timer.Stop()

	// 初次死线
	conn.SetReadDeadline(time.Now().Add(holdWait))
	// 设置收pong触发器
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(holdWait))
		timer.Reset(pingSpit)
		return nil
	})

	s.ctx, s.cancel = context.WithCancel(context.Background())
	defer s.cancel()
	go s.setOnline(true)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("连接异常:%v", err)
			break
		}
		conn.SetReadDeadline(time.Now().Add(holdWait))
		timer.Reset(pingSpit)

		msg, err := MessageFromJSON(msgBytes)
		if err != nil {
			log.Printf("解析失败:%v", err)
			continue
		}
		log.Printf("%s", formatMessage(msg))
	}
}

func (s *sPeer) setOnline(re bool) {
	if re == false {
		s.status = false
		return
	}
	s.status = true
	defer s.setOnline(false)

	for {
		select {
		case <-s.timer.C:
			err := s.conn.WriteControl(
				websocket.PingMessage,
				nil,
				time.Now().Add(writeWait),
			)
			if err != nil {
				return
			}
			s.timer.Reset(pingSpit)
		case <-s.ctx.Done():
			return
		}
	}
}
