package router

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Zoomish/portfolio-backend/internal/bot/callback"
	"github.com/Zoomish/portfolio-backend/internal/bot/command"
	"github.com/Zoomish/portfolio-backend/internal/bot/screen"
	"github.com/Zoomish/portfolio-backend/internal/bot/session"
	"github.com/Zoomish/portfolio-backend/internal/bot/stats"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
	"github.com/Zoomish/portfolio-backend/internal/domain"
	"github.com/Zoomish/portfolio-backend/internal/telegram"
)

type Users interface {
	Upsert(ctx context.Context, telegramID int64, name, username string, admin bool) (domain.User, error)
}

type Commands interface {
	Handle(ctx context.Context, t *turn.Turn, name string) error
}

type Callbacks interface {
	Handle(ctx context.Context, t *turn.Turn, data string) error
}

type Inputs interface {
	Handle(ctx context.Context, t *turn.Turn, text string) error
}

type Router struct {
	Users     Users
	Sessions  session.Store
	Screen    *screen.Screen
	Commands  Commands
	Callbacks Callbacks
	Input     Inputs
	AdminID   int64
}

func (r *Router) Handle(ctx context.Context, u telegram.Update) error {
	inc, ok := incoming(u)
	if !ok || (inc.ChatType != "" && inc.ChatType != "private") {
		return nil
	}
	name := strings.TrimSpace(inc.From.FirstName)
	if name == "" {
		name = "—"
	}
	user, err := r.Users.Upsert(ctx, inc.From.ID, name, strings.TrimSpace(inc.From.Username), stats.Allowed(inc.From.ID, r.AdminID))
	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}
	state, err := r.Sessions.Load(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("load session: %w", err)
	}
	t := &turn.Turn{
		UpdateID:   u.UpdateID,
		ChatID:     inc.ChatID,
		Actor:      user,
		Session:    state,
		IncomingID: inc.MessageID,
		CallbackID: inc.CallbackID,
	}
	head := ""
	if inc.Data != "" {
		head, _ = callback.Split(inc.Data)
	}
	cmd, _ := command.Name(inc.Text)
	slog.Info("bot update",
		"telegram_id", user.TelegramID,
		"update_id", u.UpdateID,
		"waiting_for", state.WaitingFor,
		"callback", head,
		"command", cmd,
	)
	if inc.Data != "" {
		if state.Waiting() {
			return r.Screen.Answer(ctx, inc.CallbackID, session.Hint(state.WaitingFor))
		}
		return r.Callbacks.Handle(ctx, t, inc.Data)
	}
	if cmd != "" {
		return r.Commands.Handle(ctx, t, cmd)
	}
	if state.Waiting() {
		return r.Input.Handle(ctx, t, inc.Text)
	}
	if strings.TrimSpace(inc.Text) == "" {
		return nil
	}
	if err := r.Screen.SendUserError(ctx, t, "Напиши /start, чтобы открыть меню."); err != nil {
		return err
	}
	return r.Screen.DeleteIncoming(ctx, t)
}

type inbound struct {
	ChatID     int64
	ChatType   string
	From       telegram.User
	Text       string
	Data       string
	MessageID  int64
	CallbackID string
}

func incoming(u telegram.Update) (inbound, bool) {
	if u.Callback != nil {
		if u.Callback.From.ID == 0 || u.Callback.ChatID == 0 {
			return inbound{}, false
		}
		return inbound{
			ChatID:     u.Callback.ChatID,
			ChatType:   u.Callback.ChatType,
			From:       u.Callback.From,
			Data:       u.Callback.Data,
			CallbackID: u.Callback.ID,
		}, true
	}
	if u.Message == nil || u.Message.From.ID == 0 || u.Message.ChatID == 0 {
		return inbound{}, false
	}
	return inbound{
		ChatID:    u.Message.ChatID,
		ChatType:  u.Message.ChatType,
		From:      u.Message.From,
		Text:      u.Message.Text,
		MessageID: u.Message.ID,
	}, true
}
