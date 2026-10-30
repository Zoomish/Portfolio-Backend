package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type Admin struct {
	token  string
	chatID int64
	http   *http.Client
}

func NewAdmin(token string, chatID int64) *Admin {
	return &Admin{
		token:  token,
		chatID: chatID,
		http:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (a *Admin) Alert(ctx context.Context, text string) {
	if a == nil || a.token == "" || a.chatID == 0 {
		slog.Error("admin alert skipped", "reason", "telegram is not configured")
		return
	}
	text = Clip(text, 500)
	body, err := json.Marshal(map[string]any{
		"chat_id": a.chatID,
		"text":    text,
	})
	if err != nil {
		slog.Error("admin alert encode", "err", err)
		return
	}
	endpoint := "https://api.telegram.org/bot" + a.token + "/sendMessage"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		slog.Error("admin alert request", "err", failure(err))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.http.Do(req)
	if err != nil {
		slog.Error("admin alert failed", "err", failure(err))
		return
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, io.LimitReader(resp.Body, 4096)); err != nil {
		slog.Error("admin alert drain", "err", err)
	}
	if resp.StatusCode >= 300 {
		slog.Error("admin alert failed", "status", resp.StatusCode)
	}
}

func failure(err error) string {
	var urlErr *url.Error
	if errors.As(err, &urlErr) && urlErr.Err != nil {
		return urlErr.Err.Error()
	}
	return "request failed"
}

func Clip(text string, limit int) string {
	if limit <= 0 || len(text) <= limit {
		return text
	}
	return text[:limit]
}
