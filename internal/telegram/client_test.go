package telegram

import (
	"errors"
	"net/url"
	"strings"
	"testing"
)

func TestFailureOmitsToken(t *testing.T) {
	err := &url.Error{
		Op:  "Post",
		URL: "https://api.telegram.org/botSECRET/sendMessage",
		Err: errors.New("timeout"),
	}
	got := failure(err)
	if strings.Contains(got, "SECRET") || strings.Contains(got, "api.telegram.org") {
		t.Fatalf("logged url: %s", got)
	}
	if got != "timeout" {
		t.Fatalf("got %s", got)
	}
}

func TestDecodeUpdates(t *testing.T) {
	raw := []byte(`[{
		"update_id": 10,
		"message": {
			"message_id": 4,
			"from": {"id": 7, "first_name": "Egor", "username": "zoomish"},
			"chat": {"id": 7, "type": "private"},
			"text": "/start"
		}
	},{
		"update_id": 11,
		"callback_query": {
			"id": "cb",
			"from": {"id": 7, "first_name": "Egor"},
			"message": {"message_id": 4, "chat": {"id": 7, "type": "private"}},
			"data": "menu-about"
		}
	}]`)
	got, err := decodeUpdates(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len %d", len(got))
	}
	if got[0].Message == nil || got[0].Message.Text != "/start" || got[0].TelegramID() != 7 {
		t.Fatalf("message %+v", got[0])
	}
	if got[1].Callback == nil || got[1].Callback.Data != "menu-about" || got[1].Callback.ChatID != 7 {
		t.Fatalf("callback %+v", got[1].Callback)
	}
}
