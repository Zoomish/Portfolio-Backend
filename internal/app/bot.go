package app

import (
	"context"
	"strings"

	"github.com/Zoomish/portfolio-backend/internal/bot/callback"
	"github.com/Zoomish/portfolio-backend/internal/bot/command"
	"github.com/Zoomish/portfolio-backend/internal/bot/input"
	"github.com/Zoomish/portfolio-backend/internal/bot/navigate"
	"github.com/Zoomish/portfolio-backend/internal/bot/poll"
	"github.com/Zoomish/portfolio-backend/internal/bot/router"
	"github.com/Zoomish/portfolio-backend/internal/bot/screen"
	"github.com/Zoomish/portfolio-backend/internal/bot/snap"
	"github.com/Zoomish/portfolio-backend/internal/bot/stats"
	"github.com/Zoomish/portfolio-backend/internal/config"
	"github.com/Zoomish/portfolio-backend/internal/content/serve"
	eventstore "github.com/Zoomish/portfolio-backend/internal/storage/event"
	sessionstore "github.com/Zoomish/portfolio-backend/internal/storage/session"
	userstore "github.com/Zoomish/portfolio-backend/internal/storage/user"
	"github.com/Zoomish/portfolio-backend/internal/telegram"
	"github.com/jackc/pgx/v5/pgxpool"
)

func wireBot(pool *pgxpool.Pool, cfg config.Config, snaps *serve.Service, alert func(context.Context, string)) starter {
	if strings.TrimSpace(cfg.TelegramToken) == "" {
		return nil
	}
	api := telegram.New(cfg.TelegramToken)
	users := userstore.New(pool)
	sessions := sessionstore.New(pool)
	events := eventstore.New(pool)
	screens := screen.New(api, sessions)
	query := stats.NewQuery(users, events)
	nav := &navigate.Nav{Screen: screens, Snaps: snap.Reader{Snaps: snaps}, Events: events, Alert: alert}
	commands := &command.Handler{
		Screen: screens, Nav: nav, Events: events, Stats: query, Alert: alert, AdminID: cfg.AdminTelegramID,
	}
	callbacks := &callback.Handler{
		Screen: screens, Nav: nav, Stats: query, Alert: alert, AdminID: cfg.AdminTelegramID,
	}
	in := &input.Handler{Screen: screens, Alert: alert, AdminID: cfg.AdminTelegramID}
	updates := &router.Router{
		Users: users, Sessions: sessions, Screen: screens,
		Commands: commands, Callbacks: callbacks, Input: in, AdminID: cfg.AdminTelegramID,
	}
	return poll.New(api.Updates, updates.Handle, api.SetCommands, alert, cfg.AlertAfter)
}
