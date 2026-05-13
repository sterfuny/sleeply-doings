package main

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Peer) UpdateInfo(name, pkg string, Battery int, Screen bool) {
	c.muPub.Lock()
	defer c.muPub.Unlock()
	if c.lastMessage == nil {
		c.lastMessage = &Message{}
	}
	if c.lastMessage.App == nil {
		c.lastMessage.App = &AppInfo{name, pkg}
	}
	if c.lastMessage.Battery == nil {
		c.lastMessage.Battery = &Battery
	}
	if c.lastMessage.Screen == nil {
		c.lastMessage.Screen = &Screen
	}
}

func (c *Peer) connect() error {
	u, err := url.Parse(c.serverURL)
	if err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
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

	log.Printf("已建立连接:%s", c.serverURL)

	// 重连补发完整状态
	msg := c.lastMessage
	if msg.App != nil || msg.Battery != nil || msg.Screen != nil {
		err := c.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Peer) Send(msg *Message) error {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()

	if conn == nil {
		c.connect()
		// return ErrNotConnected
	}

	data, err := msg.ToJSON()
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		c.connect()
		return err
	}

	// 发送成功，重置空闲计时器
	if c.timer != nil {
		c.timer.Reset(holdWait)
	}

	return nil
}

func (c *Peer) Close() {
	c.cancel()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
	}
}

