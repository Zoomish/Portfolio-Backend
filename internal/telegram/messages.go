package telegram

import (
	"context"
	"errors"
	"strings"
)

func (c *Client) SendMessage(ctx context.Context, chatID int64, text string, keys Keyboard) (int64, error) {
	payload := map[string]any{
		"chat_id":                  chatID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
		"reply_markup":             keyboard(keys),
	}
	var sent struct {
		MessageID int64 `json:"message_id"`
	}
	if err := c.call(ctx, "sendMessage", payload, &sent); err != nil {
		return 0, err
	}
	return sent.MessageID, nil
}

func (c *Client) EditMessageText(ctx context.Context, chatID, messageID int64, text string, keys Keyboard) error {
	payload := map[string]any{
		"chat_id":                  chatID,
		"message_id":               messageID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
		"reply_markup":             keyboard(keys),
	}
	return c.call(ctx, "editMessageText", payload, nil)
}

func (c *Client) DeleteMessage(ctx context.Context, chatID, messageID int64) error {
	payload := map[string]any{
		"chat_id":    chatID,
		"message_id": messageID,
	}
	return c.call(ctx, "deleteMessage", payload, nil)
}

func (c *Client) AnswerCallback(ctx context.Context, callbackID, text string) error {
	payload := map[string]any{"callback_query_id": callbackID}
	if text != "" {
		payload["text"] = clip(text, 180)
	}
	return c.call(ctx, "answerCallbackQuery", payload, nil)
}

func (c *Client) SetCommands(ctx context.Context) error {
	payload := map[string]any{
		"commands": []map[string]string{
			{"command": "start", "description": "Начать"},
			{"command": "settings", "description": "Меню"},
			{"command": "error", "description": "Сообщить о проблеме"},
			{"command": "stats", "description": "Статистика"},
		},
	}
	return c.call(ctx, "setMyCommands", payload, nil)
}

func clip(text string, limit int) string {
	runes := []rune(text)
	if limit <= 0 || len(runes) <= limit {
		return text
	}
	return string(runes[:limit])
}

func NotModified(err error) bool {
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr == nil {
		return false
	}
	return strings.Contains(strings.ToLower(apiErr.Description), "message is not modified")
}
