package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/config"
	"github.com/Zoomish/portfolio-backend/internal/content/github"
	"github.com/Zoomish/portfolio-backend/internal/content/linkedin"
	"github.com/Zoomish/portfolio-backend/internal/content/refresh"
	"github.com/Zoomish/portfolio-backend/internal/content/serve"
	"github.com/Zoomish/portfolio-backend/internal/httpapi"
	kajob "github.com/Zoomish/portfolio-backend/internal/keepalive"
	"github.com/Zoomish/portfolio-backend/internal/notify"
	"github.com/Zoomish/portfolio-backend/internal/postgres"
	kastore "github.com/Zoomish/portfolio-backend/internal/storage/keepalive"
	"github.com/Zoomish/portfolio-backend/internal/storage/snapshot"
	"github.com/jackc/pgx/v5/pgxpool"
)

type runtime struct {
	pool      *pgxpool.Pool
	server    *http.Server
	refresh   *refresh.Loop
	keepalive *kajob.Loop
}

func wire(ctx context.Context, cfg config.Config) (*runtime, error) {
	pool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return nil, err
	}
	if err := postgres.Migrate(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}

	snaps := snapshot.New(pool)
	beats := kastore.New(pool)
	repos := github.NewClient(cfg.GitHubToken)
	profiles := linkedin.NewClient(cfg.RapidAPIKey, cfg.RapidAPIHost, cfg.LinkedInUser)
	admin := notify.NewAdmin(cfg.TelegramToken, cfg.AdminTelegramID)
	cycle := refresh.NewCycle(repos, profiles, snaps)
	due := func(ctx context.Context) bool {
		needed, err := refresh.NeedsRefresh(ctx, snaps, time.Now(), cfg.RefreshEvery)
		if err != nil {
			slog.Error("refresh check", "err", err)
			return true
		}
		return needed
	}
	updates := refresh.NewLoop(cycle.Run, cfg.RefreshEvery, cfg.AlertAfter, admin.Alert, due)
	contentAPI := serve.New(snaps, updates, serve.SystemClock{}, cfg.RefreshEvery)
	keeper := kajob.NewLoop(beats, cfg.KeepaliveEvery, cfg.KeepaliveEnabled, cfg.AlertAfter, admin.Alert)

	return &runtime{
		pool: pool,
		server: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           httpapi.New(contentAPI, contentAPI, beats),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		refresh:   updates,
		keepalive: keeper,
	}, nil
}

func (rt *runtime) close() {
	rt.pool.Close()
}
