package api

import (
	"net/http"
	"strconv"

	. "sleeply-alive/internal/models"
	"sleeply-alive/internal/ctrl"
)

type PushHandler struct{}

func (h *PushHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "非GET/POST请求", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	App := &AppInfo{}
	var Battery *int
	var Screen *bool

	if v := q.Get("app_name"); v != "" {
		App.Name = v
	}

	if v := q.Get("app_pkg"); v != "" {
		App.Pkg = v
	}

	if v := q.Get("battery"); v != "" {
		if b, err := strconv.Atoi(v); err == nil {
			Battery = &b
		}
	}

	if v := q.Get("screen"); v != "" {
		if s, err := strconv.ParseBool(v); err == nil {
			Screen = &s
		}
	}

	ctrl.UpdateInfo(App, Battery, Screen)
	w.WriteHeader(http.StatusOK) //return 200
}
