package refresh

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/notify"
)

const backgroundTimeout = 2 * time.Minute

type Loop struct {
	run    func(context.Context) error
	every  time.Duration
	alert  func(context.Context, string)
	due    func(context.Context) bool
	streak *notify.Streak

	mu      sync.Mutex
	base    context.Context
	running bool
	done    chan struct{}
	last    error
}

func NewLoop(run func(context.Context) error, every time.Duration, alertAfter int, alert func(context.Context, string), due func(context.Context) bool) *Loop {
	return &Loop{
		run:    run,
		every:  every,
		alert:  alert,
		due:    due,
		streak: notify.NewStreak(alertAfter),
	}
}

func (l *Loop) Bind(ctx context.Context) {
	l.mu.Lock()
	l.base = ctx
	l.mu.Unlock()
}

func (l *Loop) Start(ctx context.Context) {
	l.Bind(ctx)
	if l.due != nil && l.due(ctx) {
		l.Trigger()
	}
	ticker := time.NewTicker(l.every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.Trigger()
		}
	}
}

func (l *Loop) Trigger() {
	l.mu.Lock()
	if l.running {
		l.mu.Unlock()
		return
	}
	parent := l.base
	l.running = true
	l.done = make(chan struct{})
	l.mu.Unlock()
	if parent == nil {
		parent = context.Background()
	}

	go func() {
		ctx, cancel := context.WithTimeout(parent, backgroundTimeout)
		defer cancel()
		l.finish(l.run(ctx))
	}()
}

func (l *Loop) RunOnce(ctx context.Context) error {
	l.mu.Lock()
	if l.running {
		done := l.done
		l.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			l.mu.Lock()
			err := l.last
			l.mu.Unlock()
			return err
		}
	}
	l.running = true
	l.done = make(chan struct{})
	l.mu.Unlock()

	err := l.run(ctx)
	l.finish(err)
	return err
}

func (l *Loop) finish(err error) {
	l.observe(err)
	l.mu.Lock()
	l.last = err
	l.running = false
	close(l.done)
	l.mu.Unlock()
}

func (l *Loop) observe(err error) {
	if err == nil {
		l.streak.OK()
		slog.Info("content refresh ok")
		return
	}
	if errors.Is(err, context.Canceled) {
		return
	}
	slog.Error("content refresh", "err", err)
	if l.streak.Fail() && l.alert != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		l.alert(ctx, "content refresh failed: "+notify.Clip(err.Error(), 400))
	}
}
