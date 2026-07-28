package api

import (
	"net/http"
	"strconv"
	"sleeply-alive/internal/ctrl"
	// . "sleeply-alive/internal/models"
)

func autoConvert(s string) interface{} {
	// 尝试 int
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	// 尝试 float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	// 尝试 bool
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	// default
	return s
}

type PushHandler struct{}

func (h *PushHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "非GET/POST请求", http.StatusMethodNotAllowed)
		return
	}

	params := make(map[string]any)
	for k, v := range r.URL.Query() {
		if len(v) == 1 {
			params[k] = autoConvert(v[0])
		} else {
			params[k] = v
		}
	}

	ctrl.UpdateInfo(params)
	w.WriteHeader(http.StatusOK) //return 200
}
