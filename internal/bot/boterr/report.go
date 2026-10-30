package boterr

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Zoomish/portfolio-backend/internal/bot/screen"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
)

func Report(ctx context.Context, alert func(context.Context, string), scr *screen.Screen, t *turn.Turn, err error) error {
	slog.Error("bot", "err", err, "telegram_id", t.Actor.TelegramID, "update_id", t.UpdateID, "waiting_for", t.Session.WaitingFor)
	if alert != nil {
		alert(ctx, fmt.Sprintf("bot: %v", err))
	}
	if sendErr := scr.SendUserError(ctx, t, "Что-то пошло не так. Попробуй позже."); sendErr != nil {
		return fmt.Errorf("send user error: %w", sendErr)
	}
	return nil
}
