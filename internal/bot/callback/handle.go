package callback

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/bot/boterr"
	"github.com/Zoomish/portfolio-backend/internal/bot/card"
	"github.com/Zoomish/portfolio-backend/internal/bot/navigate"
	"github.com/Zoomish/portfolio-backend/internal/bot/screen"
	"github.com/Zoomish/portfolio-backend/internal/bot/stats"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
)

type Handler struct {
	Screen  *screen.Screen
	Nav     *navigate.Nav
	Stats   stats.Query
	Alert   func(context.Context, string)
	AdminID int64
	Now     func() time.Time
}

func (h *Handler) Handle(ctx context.Context, t *turn.Turn, data string) error {
	head, rest := Split(data)
	switch head {
	case "menu":
		return h.menu(ctx, t, rest)
	case "stats":
		return h.stats(ctx, t, rest)
	default:
		return h.Screen.Answer(ctx, t.CallbackID, "Неизвестная кнопка")
	}
}

func (h *Handler) menu(ctx context.Context, t *turn.Turn, rest string) error {
	kind, ok := card.Kind(rest)
	if !ok {
		return h.Screen.Answer(ctx, t.CallbackID, "Неизвестная кнопка")
	}
	h.ack(ctx, t)
	return h.Nav.Open(ctx, t, navigate.ModeEdit, kind)
}

func (h *Handler) stats(ctx context.Context, t *turn.Turn, rest string) error {
	if !stats.Allowed(t.Actor.TelegramID, h.AdminID) {
		return h.Screen.Answer(ctx, t.CallbackID, "Недоступно")
	}
	switch rest {
	case "recent":
		h.ack(ctx, t)
		people, err := h.Stats.Recent(ctx, stats.RecentLimit)
		if err != nil {
			return boterr.Report(ctx, h.Alert, h.Screen, t, fmt.Errorf("stats recent: %w", err))
		}
		return h.Screen.EditFunctional(ctx, t, stats.FormatRecent(people), stats.RecentKeys())
	case "summary":
		h.ack(ctx, t)
		summary, err := h.Stats.Summary(ctx, h.now())
		if err != nil {
			return boterr.Report(ctx, h.Alert, h.Screen, t, fmt.Errorf("stats: %w", err))
		}
		return h.Screen.EditFunctional(ctx, t, stats.FormatSummary(summary), stats.SummaryKeys())
	default:
		return h.Screen.Answer(ctx, t.CallbackID, "Неизвестная кнопка")
	}
}

func (h *Handler) ack(ctx context.Context, t *turn.Turn) {
	if err := h.Screen.Answer(ctx, t.CallbackID, ""); err != nil {
		slog.Warn("answer callback", "err", err, "telegram_id", t.Actor.TelegramID)
	}
}

func (h *Handler) now() time.Time {
	if h.Now != nil {
		return h.Now()
	}
	return time.Now()
}
