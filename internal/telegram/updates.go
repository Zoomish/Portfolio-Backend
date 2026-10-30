package telegram

import (
	"context"
	"encoding/json"
	"fmt"
)

type User struct {
	ID        int64
	FirstName string
	Username  string
}

type Incoming struct {
	ID       int64
	ChatID   int64
	ChatType string
	From     User
	Text     string
}

type Callback struct {
	ID        string
	From      User
	ChatID    int64
	ChatType  string
	MessageID int64
	Data      string
}

type Update struct {
	UpdateID int64
	Message  *Incoming
	Callback *Callback
}

func (u Update) TelegramID() int64 {
	if u.Callback != nil {
		return u.Callback.From.ID
	}
	if u.Message != nil {
		return u.Message.From.ID
	}
	return 0
}

func (c *Client) Updates(ctx context.Context, offset int64) ([]Update, error) {
	payload := map[string]any{
		"offset":          offset,
		"timeout":         25,
		"allowed_updates": []string{"message", "callback_query"},
	}
	var raw json.RawMessage
	if err := c.call(ctx, "getUpdates", payload, &raw); err != nil {
		return nil, err
	}
	return decodeUpdates(raw)
}

func decodeUpdates(raw []byte) ([]Update, error) {
	var batch []rawUpdate
	if len(raw) == 0 {
		return nil, nil
	}
	if err := json.Unmarshal(raw, &batch); err != nil {
		return nil, fmt.Errorf("updates decode: %w", err)
	}
	out := make([]Update, 0, len(batch))
	for _, item := range batch {
		out = append(out, item.update())
	}
	return out, nil
}

type rawUpdate struct {
	UpdateID int64        `json:"update_id"`
	Message  *rawMessage  `json:"message"`
	Callback *rawCallback `json:"callback_query"`
}

type rawMessage struct {
	MessageID int64    `json:"message_id"`
	From      *rawUser `json:"from"`
	Chat      rawChat  `json:"chat"`
	Text      string   `json:"text"`
}

type rawChat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type rawUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type rawCallback struct {
	ID      string      `json:"id"`
	From    rawUser     `json:"from"`
	Message *rawMessage `json:"message"`
	Data    string      `json:"data"`
}

func (u rawUpdate) update() Update {
	out := Update{UpdateID: u.UpdateID}
	if u.Message != nil && u.Message.From != nil {
		out.Message = &Incoming{
			ID:       u.Message.MessageID,
			ChatID:   u.Message.Chat.ID,
			ChatType: u.Message.Chat.Type,
			From:     u.Message.From.user(),
			Text:     u.Message.Text,
		}
	}
	if u.Callback != nil {
		cb := &Callback{
			ID:   u.Callback.ID,
			From: u.Callback.From.user(),
			Data: u.Callback.Data,
		}
		if u.Callback.Message != nil {
			cb.ChatID = u.Callback.Message.Chat.ID
			cb.ChatType = u.Callback.Message.Chat.Type
			cb.MessageID = u.Callback.Message.MessageID
		}
		out.Callback = cb
	}
	return out
}

func (u rawUser) user() User {
	return User{ID: u.ID, FirstName: u.FirstName, Username: u.Username}
}
