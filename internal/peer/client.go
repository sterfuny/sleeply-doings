package peer

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"
	"sync"

	"github.com/gorilla/websocket"
	cfg "sleeply-alive/internal/config"
	. "sleeply-alive/internal/models"
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

	registerMsg := fmt.Sprintf(`{"id":"%s"}`, cfg.Get().ID)
	conn.WriteMessage(websocket.TextMessage, []byte(registerMsg))
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

func (c *CPeer) Send(msg any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn := c.conn

	if conn == nil {
		c.Connect()
	}

	var data []byte

	switch v := msg.(type) {
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

func (c *CPeer) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
	}
}
