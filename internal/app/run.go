package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/config"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cfg config.Config) error {
	SetupLog(cfg.Env, "")

	rt, err := wire(ctx, cfg)
	if err != nil {
		return err
	}
	defer rt.close()

	group, ctx := errgroup.WithContext(ctx)
	rt.refresh.Bind(ctx)
	group.Go(func() error {
		slog.Info("http listen", "addr", rt.server.Addr)
		err := rt.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("http listen: %w", err)
		}
		return nil
	})
	group.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := rt.server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http shutdown: %w", err)
		}
		return nil
	})
	group.Go(func() error {
		rt.refresh.Start(ctx)
		return nil
	})
	group.Go(func() error {
		rt.keepalive.Start(ctx)
		return nil
	})
	group.Go(func() error {
		if rt.bot == nil {
			<-ctx.Done()
			return nil
		}
		rt.bot.Start(ctx)
		return nil
	})
	if err := group.Wait(); err != nil {
		return fmt.Errorf("run: %w", err)
	}
	return nil
}
