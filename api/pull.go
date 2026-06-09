package api

import (
	"net/http"

	. "sleeply-alive/internal/models"
	"sleeply-alive/internal/ctrl"
)

type PullHandler struct{}

func (h *PullHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "非GET/POST请求", http.StatusMethodNotAllowed)
		return
	}

	var dev *Device
	q := r.URL.Query()
	// path := r.URL.Path

	if v := q.Get("uuid"); v != "" {
		dev = ctrl.FindDev(v)
		if dev == nil {
			http.Error(w, "nil", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK) //return 200
	body, err := dev.ToJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
