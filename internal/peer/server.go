package peer

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	. "sleeply-alive/internal/models"
)

type SPeer struct {
	device *Device
	conn   *websocket.Conn
	timer  *time.Timer

	ctx    context.Context
	cancel context.CancelFunc
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *SPeer)findId(data []byte) error {
	var tmp struct {
		ID string `json:"id"`
	}

	json.Unmarshal(data, &tmp)
	if _, err := uuid.Parse(tmp.ID); err != nil {
		return err
	}
	id := tmp.ID
	s.device = Find(id)
	return nil
}

func handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	ws := &SPeer{conn: conn}
	go ws.handleWSClient()
}

func (s *SPeer) handleWSClient() {
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

	_, msgBytes, _ := conn.ReadMessage()
	if err := s.findId(msgBytes); err != nil {
		log.Printf("未知设备")
	}
	go s.setOnline(true)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("连接异常:%v", err)
			break
		}
		conn.SetReadDeadline(time.Now().Add(holdWait))
		timer.Reset(pingSpit)

		s.device.Lastmsg, err = MessageFromJSON(msgBytes)
		if err != nil {
			log.Printf("解析失败:%v", err)
			continue
		}
		log.Printf("%s", /*formatMessage(s.device.Lastmsg)*/)
	}
}

func (s *SPeer) setOnline(re bool) {
	if re == false {
		s.device.Status = false
		return
	}
	s.device.Status = true
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

func RegisterHandle(addr string) error {
	http.HandleFunc("/ws", handleConn)
	log.Printf("服务启动:%s", addr)
	return http.ListenAndServe(addr, nil)
}
