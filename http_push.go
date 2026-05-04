package main

import (
	"log"
	"net/http"
	"strconv"
)

type PushHandler struct {
	client *Client
}

func (h *PushHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "仅支持GET/POST", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	var msg Message

	if v := q.Get("app_name"); v != "" {
		if msg.App == nil {
			msg.App = &AppInfo{}
		}
		msg.App.Name = v
	}

	if v := q.Get("app_pkg"); v != "" {
		if msg.App == nil {
			msg.App = &AppInfo{}
		}
		msg.App.Pkg = v
	}

	if v := q.Get("battery"); v != "" {
		if b, err := strconv.Atoi(v); err == nil {
			msg.Battery = &b
		}
	}

	if v := q.Get("screen"); v != "" {
		if s, err := strconv.ParseBool(v); err == nil {
			msg.Screen = &s
		}
	}
	
	// 更新缓存
	if msg.App != nil {
		h.client.UpdateApp(msg.App.Name, msg.App.Pkg)
	}
	if msg.Battery != nil {
		h.client.UpdateBattery(*msg.Battery)
	}
	if msg.Screen != nil {
		h.client.UpdateScreen(*msg.Screen)
	}

	//resp, err := h.client.Send(&msg)
	if err := h.client.Send(&msg); err != nil {
    log.Printf("WebSocket推送失败: %v", err)
    http.Error(w, "推送失败: "+err.Error(), http.StatusServiceUnavailable)
    return
	}
	
	w.WriteHeader(http.StatusOK)  //return 200
	log.Printf("HTTP推送成功: %s", formatMessage(&msg))
}

func startHTTPPush(client *Client, addr string) error {
	handler := &PushHandler{client: client}
	log.Printf("HTTP接口启动在 %s", addr)
	return http.ListenAndServe(addr, handler)
}