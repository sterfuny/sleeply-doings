package api

import (
	"log"
	"net/http"
)

func StartHTTPPush(addr string) error {
	log.Printf("push接口启动在%s", addr)
	return http.ListenAndServe(addr, &PushHandler{})
}
func StartHTTPPull(addr string) error {
	log.Printf("pull接口启动在%s", addr)
	return http.ListenAndServe(addr, &PullHandler{})
}
