package api

import (
	"net/http"

	"sleeply-alive/internal/ctrl"
	. "sleeply-alive/internal/models"
	"sleeply-alive/internal/peer"
)

type PullHandler struct{}

func (h *PullHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/ws" {
		h.handleConn(w, r)
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "非GET/POST请求", http.StatusMethodNotAllowed)
		return
	}

	var dev *Device
	q := r.URL.Query()

	if v := q.Get("uuid"); v != "" {
		dev = ctrl.FindDev(v)
		if dev == nil {
			http.Error(w, "nil", http.StatusBadRequest)
			return
		}
	}

	body, err := dev.ToJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (h *PullHandler) handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	peer.PeerInit(conn)
}
