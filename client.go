package main

import (
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

	c.mu.Unlock()
	log.Printf("已建立连接:%s", c.URL)

	if msg.App != nil || msg.Battery != nil || msg.Screen != nil {
		err := c.Send(msg)
		if err != nil {
			return err
		}
	}
	for {
		_, _, err := conn.ReadMessage()

		if err != nil {
			return err
		}
		conn.SetReadDeadline(time.Now().Add(holdWait))
	}
}

func (c *cPeer) Send(msg *Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn := c.conn

	if conn == nil {
		c.connect()
	}

	data, err := msg.ToJSON()
	// log.Printf("jsonsend")
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}
	conn.SetReadDeadline(time.Now().Add(holdWait))
	return nil
}

func (c *cPeer) Close() {
	c.cancel()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
	}
}
