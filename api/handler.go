package api

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartHTTPPush(addr string) error {
	log.Printf("push接口启动在%s", addr)
	return http.ListenAndServe(addr, &PushHandler{})
}
func StartHTTPPull(addr string) error {
	log.Printf("pull接口启动在%s", addr)
	return http.ListenAndServe(addr, &PullHandler{})
}
