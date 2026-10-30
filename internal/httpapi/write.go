package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Zoomish/portfolio-backend/internal/content"
)

type errBody struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("write json", "err", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, content.ErrNotReady) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		writeJSON(w, http.StatusServiceUnavailable, errBody{Error: "content is not ready"})
		return
	}
	slog.Error("request failed", "err", err)
	writeJSON(w, http.StatusInternalServerError, errBody{Error: "something went wrong"})
}
