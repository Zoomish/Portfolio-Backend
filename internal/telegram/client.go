package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	token string
	base  string
	http  *http.Client
}

func New(token string) *Client {
	return &Client{
		token: token,
		base:  "https://api.telegram.org/bot",
		http:  &http.Client{Timeout: 40 * time.Second},
	}
}

type Button struct {
	Text string
	Data string
}

type Keyboard [][]Button

type Error struct {
	Method      string
	Code        int
	Description string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Method, e.Description)
}

func (c *Client) call(ctx context.Context, method string, payload any, result any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("%s encode: %w", method, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+c.token+"/"+method, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%s: %s", method, failure(err))
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return context.Canceled
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return context.DeadlineExceeded
		}
		return fmt.Errorf("%s: %s", method, failure(err))
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("%s read: %w", method, err)
	}
	var envelope struct {
		OK          bool            `json:"ok"`
		Description string          `json:"description"`
		ErrorCode   int             `json:"error_code"`
		Result      json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return fmt.Errorf("%s decode: %w", method, err)
	}
	if !envelope.OK {
		return &Error{Method: method, Code: envelope.ErrorCode, Description: envelope.Description}
	}
	if result == nil || len(envelope.Result) == 0 {
		return nil
	}
	if err := json.Unmarshal(envelope.Result, result); err != nil {
		return fmt.Errorf("%s result: %w", method, err)
	}
	return nil
}

func failure(err error) string {
	var urlErr *url.Error
	if errors.As(err, &urlErr) && urlErr.Err != nil {
		return urlErr.Err.Error()
	}
	return "request failed"
}

type markup struct {
	InlineKeyboard [][]inlineButton `json:"inline_keyboard"`
}

type inlineButton struct {
	Text string `json:"text"`
	Data string `json:"callback_data"`
}

func keyboard(keys Keyboard) markup {
	rows := make([][]inlineButton, 0, len(keys))
	for _, row := range keys {
		buttons := make([]inlineButton, 0, len(row))
		for _, button := range row {
			buttons = append(buttons, inlineButton{Text: button.Text, Data: button.Data})
		}
		rows = append(rows, buttons)
	}
	return markup{InlineKeyboard: rows}
}
