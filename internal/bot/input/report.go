package input

import (
	"context"
	"fmt"
	"strings"

	"github.com/Zoomish/portfolio-backend/internal/bot/screen"
	"github.com/Zoomish/portfolio-backend/internal/bot/session"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
)

type Handler struct {
	Screen  *screen.Screen
	Alert   func(context.Context, string)
	AdminID int64
}

func (h *Handler) Handle(ctx context.Context, t *turn.Turn, text string) error {
	if t.Session.WaitingFor != session.WaitErrorReport {
		if err := h.Screen.SendUserError(ctx, t, session.Hint(t.Session.WaitingFor)); err != nil {
			return err
		}
		return h.Screen.DeleteIncoming(ctx, t)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		if err := h.Screen.SendUserError(ctx, t, "Пришли текст проблемы."); err != nil {
			return err
		}
		return h.Screen.DeleteIncoming(ctx, t)
	}
	if h.AdminID == 0 {
		if err := h.Screen.EditFunctional(ctx, t, "Сейчас некому передать сообщение.", nil); err != nil {
			return err
		}
		t.Session.WaitingFor = ""
		return h.Screen.ClearUserNoise(ctx, t)
	}
	if h.Alert != nil {
		h.Alert(ctx, report(t, text))
	}
	if err := h.Screen.EditFunctional(ctx, t, "Передал. Спасибо.", nil); err != nil {
		return err
	}
	t.Session.WaitingFor = ""
	return h.Screen.ClearUserNoise(ctx, t)
}

func report(t *turn.Turn, text string) string {
	line := fmt.Sprintf("Проблема от %s", t.Actor.Name)
	if t.Actor.Username != "" {
		line += " @" + t.Actor.Username
	}
	return fmt.Sprintf("%s id %d\n%s", line, t.Actor.TelegramID, text)
}
