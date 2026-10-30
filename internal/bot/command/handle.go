package command

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/bot/card"
	"github.com/Zoomish/portfolio-backend/internal/bot/navigate"
	"github.com/Zoomish/portfolio-backend/internal/bot/screen"
	"github.com/Zoomish/portfolio-backend/internal/bot/session"
	"github.com/Zoomish/portfolio-backend/internal/bot/stats"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
	"github.com/Zoomish/portfolio-backend/internal/domain"
)

type Events interface {
	Add(ctx context.Context, userID int64, kind string) error
}

type Handler struct {
	Screen  *screen.Screen
	Nav     *navigate.Nav
	Events  Events
	Stats   stats.Query
	Alert   func(context.Context, string)
	AdminID int64
	Now     func() time.Time
}

func (h *Handler) Handle(ctx context.Context, t *turn.Turn, name string) error {
	if t.Session.Waiting() {
		if err := h.Screen.SendUserError(ctx, t, session.Hint(t.Session.WaitingFor)); err != nil {
			return err
		}
		return h.Screen.DeleteIncoming(ctx, t)
	}
	var err error
	switch name {
	case "start":
		err = h.start(ctx, t)
	case "settings":
		err = h.settings(ctx, t)
	case "error":
		err = h.report(ctx, t)
	case "stats":
		err = h.statsCmd(ctx, t)
	default:
		err = h.Screen.SendUserError(ctx, t, "Неизвестная команда. Открой /start.")
	}
	if err != nil {
		return err
	}
	return h.Screen.DeleteIncoming(ctx, t)
}

func (h *Handler) start(ctx context.Context, t *turn.Turn) error {
	if err := h.Nav.Open(ctx, t, navigate.ModeCommand, card.ScreenHome); err != nil {
		return err
	}
	if err := h.Events.Add(ctx, t.Actor.ID, domain.EventStart); err != nil {
		slog.Error("record event", "err", err, "telegram_id", t.Actor.TelegramID, "update_id", t.UpdateID)
		if h.Alert != nil {
			h.Alert(ctx, fmt.Sprintf("bot event: %v", err))
		}
	}
	return nil
}

func (h *Handler) settings(ctx context.Context, t *turn.Turn) error {
	return h.Nav.Open(ctx, t, navigate.ModeCommand, card.ScreenHome)
}

func Name(text string) (string, bool) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "/") {
		return "", false
	}
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return "", false
	}
	name, _, _ := strings.Cut(strings.TrimPrefix(fields[0], "/"), "@")
	name = strings.ToLower(name)
	if name == "" {
		return "", false
	}
	return name, true
}
