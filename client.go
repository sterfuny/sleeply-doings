package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func (c *cPeer) connect() error {
	u, err := url.Parse(c.URL)
	if err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close()
	}
	c.conn = conn

	// 初次死线
	conn.SetReadDeadline(time.Now().Add(holdWait))
	// 设置收ping触发器
	conn.SetPingHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(holdWait))
		return conn.WriteControl(
			websocket.PongMessage,
			nil,
			time.Now().Add(writeWait),
		)

	})

	registerMsg := fmt.Sprintf(`{"id":"%s"}`, Get().ID)
	conn.WriteMessage(websocket.TextMessage, []byte(registerMsg))
	log.Printf("已建立连接:%s", c.URL)
	c.mu.Unlock()
	if err := c.Send(msg); err != nil {
		return err
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		conn.SetReadDeadline(time.Now().Add(holdWait))
	}
}

func (c *cPeer) Send(msg any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn := c.conn

	if conn == nil {
		c.connect()
	}

	var data []byte

	switch v:=msg.(type) {
	case *Message:
		data, _ = v.ToJSON()
	case string:
		data = []byte(v)
	default:
		return errors.New("unkown msg type")
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}
	conn.SetReadDeadline(time.Now().Add(holdWait))
	return nil
}

func (c *cPeer) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
	}
}
