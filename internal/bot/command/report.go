package command

import (
	"context"

	"github.com/Zoomish/portfolio-backend/internal/bot/session"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
)

func (h *Handler) report(ctx context.Context, t *turn.Turn) error {
	if err := h.Screen.SendCommand(ctx, t, "Опиши проблему одним сообщением. Я передам её владельцу.", nil); err != nil {
		return err
	}
	t.Session.WaitingFor = session.WaitErrorReport
	return h.Screen.Save(ctx, t)
}
