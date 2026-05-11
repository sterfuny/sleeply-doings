// client.go
package main

import (
	"context"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	serverURL string
	conn      *websocket.Conn
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	reconnect chan struct{}

	//缓存完整状态
	lastMessage *Message
	muPub       sync.RWMutex

	//心跳计时器
	idleTimer *time.Timer
}

func (c *Client) UpdateInfo(name, pkg string, Battery int, Screen bool) {
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

func (c *Client) GetLastMessage() *Message {
	c.muPub.RLock()
	defer c.muPub.RUnlock()
	if c.lastMessage == nil {
		return &Message{}
	}
	copy := *c.lastMessage
	return &copy
}

func NewClient(serverURL string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		serverURL: serverURL,
		ctx:       ctx,
		cancel:    cancel,
		reconnect: make(chan struct{}, 1),
	}

	if err := c.connect(); err != nil {
		log.Printf("连接失败:%v", err)
	}

	//go c.reconnectLoop()

	return c
}

func (c *Client) connect() error {
	u, err := url.Parse(c.serverURL)
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
			[]byte{},
			time.Now().Add(writeWait),
		)
	})

	log.Printf("已建立连接:%s", c.serverURL)

	// 重连补发完整状态
	lastMsg := c.GetLastMessage()
	if lastMsg.App != nil || lastMsg.Battery != nil || lastMsg.Screen != nil {
		data, _ := lastMsg.ToJSON()
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("补发失败: %v", err)
		} else {
			log.Printf("重连补发: %s", formatMessage(lastMsg))
		}
	}

	// 启动读循环和空闲心跳
	go c.readLoop()
	go c.idleHeartbeat(holdWait)

	return nil
}

func (c *Client) idleHeartbeat(holdWait time.Duration) {
	for {
		select {
		case <-c.idleTimer.C:
			c.mu.Lock()
			conn := c.conn
			c.mu.Unlock()
			if conn != nil {
				if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
					return
				}
			}
			c.idleTimer.Reset(holdWait)
		case <-c.ctx.Done():
			return
		}
	}
}
/*
func (c *Client) reconnectLoop() {
	retry := 0
	maxBackoff := 60 * time.Second

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.reconnect:
			retry++

			backoff := time.Duration(math.Min(
				float64(time.Second)*math.Pow(2, float64(retry)),
				float64(maxBackoff),
			))

			log.Printf("等待%v %v...", retry, backoff)

			select {
			case <-c.ctx.Done():
				return
			case <-time.After(backoff):
			}

			if err := c.connect(); err != nil {
				log.Printf("重连失败: %v", err)
				c.triggerReconnect()
			} else {
				log.Println("重连成功")
				retry = 0
				select {
				case <-c.reconnect:
				default:
				}
			}
		}
	}
}
*/
func (c *Client) triggerReconnect() {
	select {
	case c.reconnect <- struct{}{}:
	default:
	}
}

func (c *Client) Send(msg *Message) error {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()

	if conn == nil {
		c.triggerReconnect()
		// return ErrNotConnected
	}

	data, err := msg.ToJSON()
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		c.triggerReconnect()
		return err
	}

	// 发送成功，重置空闲计时器
	if c.idleTimer != nil {
		c.idleTimer.Reset(25 * time.Second)
	}

	return nil
}

func (c *Client) readLoop() {
	for {
		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()

		if conn == nil {
			return
		}

		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取失败: %v", err)
			c.triggerReconnect()
			return
		}
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		// 收到消息，重置空闲计时器
		if c.idleTimer != nil {
			c.idleTimer.Reset(25 * time.Second)
		}
		log.Printf("收到: %s", string(msgBytes))
	}
}

func (c *Client) Close() {
	c.cancel()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
	}
}

