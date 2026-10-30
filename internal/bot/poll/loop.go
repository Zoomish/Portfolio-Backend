package poll

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/notify"
	"github.com/Zoomish/portfolio-backend/internal/telegram"
)

type Loop struct {
	updates func(context.Context, int64) ([]telegram.Update, error)
	handle  func(context.Context, telegram.Update) error
	setup   func(context.Context) error
	alert   func(context.Context, string)
	streak  *notify.Streak
}

func New(
	updates func(context.Context, int64) ([]telegram.Update, error),
	handle func(context.Context, telegram.Update) error,
	setup func(context.Context) error,
	alert func(context.Context, string),
	alertAfter int,
) *Loop {
	return &Loop{
		updates: updates,
		handle:  handle,
		setup:   setup,
		alert:   alert,
		streak:  notify.NewStreak(alertAfter),
	}
}

func (l *Loop) Start(ctx context.Context) {
	slog.Info("bot poll")
	if l.setup != nil {
		if err := l.setup(ctx); err != nil && ctx.Err() == nil {
			slog.Error("bot setup", "err", err)
		}
	}
	var offset int64
	for {
		if ctx.Err() != nil {
			return
		}
		batch, err := l.updates(ctx, offset)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			slog.Error("bot poll", "err", err)
			if l.streak.Fail() && l.alert != nil {
				l.alert(ctx, "bot poll failed")
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
			continue
		}
		l.streak.OK()
		for _, update := range batch {
			if update.UpdateID+1 > offset {
				offset = update.UpdateID + 1
			}
			l.dispatch(ctx, update)
		}
	}
}

func (l *Loop) dispatch(ctx context.Context, update telegram.Update) {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		slog.Error("bot panic", "update_id", update.UpdateID, "telegram_id", update.TelegramID(), "err", rec)
		if l.alert != nil {
			l.alert(ctx, fmt.Sprintf("bot panic update %d: %v", update.UpdateID, rec))
		}
	}()
	if err := l.handle(ctx, update); err != nil && ctx.Err() == nil {
		slog.Error("bot update", "update_id", update.UpdateID, "telegram_id", update.TelegramID(), "err", err)
		if l.alert != nil {
			l.alert(ctx, fmt.Sprintf("bot update %d: %v", update.UpdateID, err))
		}
	}
}
