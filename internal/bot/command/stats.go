package command

import (
	"context"
	"fmt"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/bot/boterr"
	"github.com/Zoomish/portfolio-backend/internal/bot/stats"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
)

func (h *Handler) statsCmd(ctx context.Context, t *turn.Turn) error {
	if !stats.Allowed(t.Actor.TelegramID, h.AdminID) {
		return h.Screen.SendUserError(ctx, t, "Эта команда только для владельца.")
	}
	summary, err := h.Stats.Summary(ctx, h.now())
	if err != nil {
		return boterr.Report(ctx, h.Alert, h.Screen, t, fmt.Errorf("stats: %w", err))
	}
	return h.Screen.SendCommand(ctx, t, stats.FormatSummary(summary), stats.SummaryKeys())
}

func (h *Handler) now() time.Time {
	if h.Now != nil {
		return h.Now()
	}
	return time.Now()
}
