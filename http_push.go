package main

import (
	"log"
	"net/http"
	"strconv"
)

type PushHandler struct{}

func (h *PushHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "非GET/POST请求", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	// msg := Message{App: &AppInfo{}}

	if v := q.Get("app_name"); v != "" {
		msg.App.Name = v
	}

	if v := q.Get("app_pkg"); v != "" {
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
	
	// msg.UpdateInfo(msg.App, msg.Battery, msg.Screen)	
	w.WriteHeader(http.StatusOK)  //return 200
	log.Printf("HTTP推送成功: %s", formatMessage(msg))
}

func startHTTPPush(addr string, ) error {
	log.Printf("HTTP接口启动在 %s", addr)
	return http.ListenAndServe(addr, &PushHandler{})
}
