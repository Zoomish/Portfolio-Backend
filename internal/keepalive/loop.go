package keepalive

import (
	"context"
	"log/slog"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/notify"
	store "github.com/Zoomish/portfolio-backend/internal/storage/keepalive"
)

type Loop struct {
	beats   *store.Store
	every   time.Duration
	enabled bool
	alert   func(context.Context, string)
	streak  *notify.Streak
}

func NewLoop(beats *store.Store, every time.Duration, enabled bool, alertAfter int, alert func(context.Context, string)) *Loop {
	return &Loop{
		beats:   beats,
		every:   every,
		enabled: enabled,
		alert:   alert,
		streak:  notify.NewStreak(alertAfter),
	}
}

func (l *Loop) Start(ctx context.Context) {
	if !l.enabled {
		<-ctx.Done()
		return
	}
	l.tick(ctx)
	ticker := time.NewTicker(l.every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.tick(ctx)
		}
	}
}

func (l *Loop) tick(ctx context.Context) {
	if err := l.beats.Touch(ctx); err != nil {
		slog.Error("db keepalive", "err", err)
		if l.streak.Fail() && l.alert != nil {
			alertCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			l.alert(alertCtx, "db keepalive failed: "+notify.Clip(err.Error(), 400))
		}
		return
	}
	l.streak.OK()
	slog.Info("db keepalive ok")
}
