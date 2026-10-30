package httpapi

import (
	"context"
	"log/slog"
	"net/http"
)

type readyCheck interface {
	Ready(ctx context.Context) error
}

type statusBody struct {
	Status string `json:"status"`
}

func (h *handlers) app(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, statusBody{Status: "ok"})
}

func (h *handlers) readiness(w http.ResponseWriter, r *http.Request) {
	if err := h.ready.Ready(r.Context()); err != nil {
		slog.Error("ready", "err", err)
		writeJSON(w, http.StatusServiceUnavailable, errBody{Error: "not ready"})
		return
	}
	writeJSON(w, http.StatusOK, statusBody{Status: "ok"})
}
