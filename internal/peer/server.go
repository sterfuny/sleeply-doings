package peer

import (
	"context"
	"fmt"
	"log"
	"maps"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"sleeply-alive/internal/ctrl"
	. "sleeply-alive/internal/model"
)

type SPeer struct {
	device *Device
	conn   *websocket.Conn
	timer  *time.Timer

	ctx    context.Context
	cancel context.CancelFunc
}

func SPeerInit(conn *websocket.Conn) {
	ws := &SPeer{conn: conn}
	ws.handleWSClient()
}

func (s *SPeer) find(tmp Message) error {
	if tmp.UUID == nil {
		return fmt.Errorf("未注册uuid")
	}

	id := *tmp.UUID
	s.device = ctrl.FindDev(id)

	if s.device == nil {
		ctrl.MkDev(id)
		s.device = ctrl.FindDev(id)
	}
	return nil
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

	// 握手后检查id信息
	_, msgBytes, _ := conn.ReadMessage()

	msg, err := ctrl.FromMessage(msgBytes)
	if err != nil {
		log.Println(err)
		return
	}
	if err := s.find(*msg); err != nil {
		log.Println(err)
		return
	}

	s.device.LastMsg = msg
	go s.setOnline(true)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("连接异常:%v", err)
			break
		}
		conn.SetReadDeadline(time.Now().Add(holdWait))
		timer.Reset(pingSpit)

		//处理消息
		msg, err := ctrl.FromMessage(msgBytes)
		prevMsg := s.device.LastMsg

		if prevMsg.SendKey != nil && msg.SendKey != nil {
			maps.Copy(*prevMsg.SendKey, *msg.SendKey)
		} else if msg.SendKey != nil {
			prevMsg.SendKey = msg.SendKey
		}

		if err != nil {
			log.Printf("解析失败:%v", err)
			continue
		}
		log.Printf("%s", ctrl.FormatMessage(msg))
	}
}

func (s *SPeer) setOnline(live bool) {
	if s.device == nil {
		log.Println("未知device")
		return
	}
	if !live {
		s.device.Status = false
		return
	}
	s.device.Status = true

	id := s.device.LastMsg.UUID.String()
	if strings.TrimSpace(s.device.Name) == "" {
		a := strings.Split(id, "-")[0]
		b := id[strings.LastIndex(id, "-")+1:]
		log.Printf("初次认证:%s-****-****-****-%s", a, b)
	} else {
		log.Printf("认证成功:%s", s.device.Name)
	}

	defer ctrl.SaveDB(s.device)
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
