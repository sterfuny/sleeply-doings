package peer

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	cfg "sleeply-alive/internal/config"
	// "sleeply-alive/internal/ctrl"
	. "sleeply-alive/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type CPeer struct {
	URL    string
	NewMsg chan *Message
	conn   *websocket.Conn

	mu    sync.Mutex
	muPub sync.RWMutex
}

func (c *CPeer) Connect() error {
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

	// registerMsg := fmt.Sprintf(`{"id":"%s"}`, cfg.Get().ID)
	// conn.WriteMessage(websocket.TextMessage, []byte(registerMsg))
	id, err :=uuid.Parse(cfg.Get().ID)
	if err != nil {
		return fmt.Errorf("无法创建id")
	}
	c.NewMsg <- &Message{UUID:&id}

	log.Printf("已建立连接:%s", c.URL)
	c.mu.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		conn.SetReadDeadline(time.Now().Add(holdWait))
	}
}

func (c *CPeer) Send(msg *Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn := c.conn

	if conn == nil {
		c.Connect()
	}

	var data []byte
	data, _ = msg.ToJSON()

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}
	conn.SetReadDeadline(time.Now().Add(holdWait))
	return nil
}

func (c *CPeer) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
	}
}
