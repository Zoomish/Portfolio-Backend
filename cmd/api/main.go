package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zoomish/portfolio-backend/internal/app"
	"github.com/Zoomish/portfolio-backend/internal/config"
)

func main() {
	app.SetupLog(os.Getenv("APP_ENV"), os.Getenv("NODE_ENV"))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, cfg); err != nil {
		slog.Error("stopped", "err", err)
		os.Exit(1)
	}
}
