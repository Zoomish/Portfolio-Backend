package screen

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Zoomish/portfolio-backend/internal/bot/session"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
	"github.com/Zoomish/portfolio-backend/internal/telegram"
)

type API interface {
	SendMessage(ctx context.Context, chatID int64, text string, keys telegram.Keyboard) (int64, error)
	EditMessageText(ctx context.Context, chatID, messageID int64, text string, keys telegram.Keyboard) error
	DeleteMessage(ctx context.Context, chatID, messageID int64) error
	AnswerCallback(ctx context.Context, callbackID, text string) error
}

type Screen struct {
	api      API
	sessions session.Store
}

func New(api API, sessions session.Store) *Screen {
	return &Screen{api: api, sessions: sessions}
}

func (s *Screen) SendCommand(ctx context.Context, t *turn.Turn, text string, keys telegram.Keyboard) error {
	id, err := s.api.SendMessage(ctx, t.ChatID, text, keys)
	if err != nil {
		return fmt.Errorf("send screen: %w", err)
	}
	t.Session.LastBotMessageID = id
	return s.Save(ctx, t)
}

func (s *Screen) EditFunctional(ctx context.Context, t *turn.Turn, text string, keys telegram.Keyboard) error {
	if t.Session.LastBotMessageID == 0 {
		return s.SendCommand(ctx, t, text, keys)
	}
	err := s.api.EditMessageText(ctx, t.ChatID, t.Session.LastBotMessageID, text, keys)
	if err == nil || telegram.NotModified(err) {
		return nil
	}
	slog.Warn("edit functional", "err", err, "telegram_id", t.Actor.TelegramID)
	return s.SendCommand(ctx, t, text, keys)
}

func (s *Screen) SendUserError(ctx context.Context, t *turn.Turn, text string) error {
	if t.Session.ErrorMessageID != 0 {
		s.delete(ctx, t, t.Session.ErrorMessageID)
	}
	id, err := s.api.SendMessage(ctx, t.ChatID, text, nil)
	if err != nil {
		return fmt.Errorf("send error: %w", err)
	}
	t.Session.ErrorMessageID = id
	return s.Save(ctx, t)
}

func (s *Screen) Answer(ctx context.Context, callbackID, text string) error {
	if callbackID == "" {
		return nil
	}
	if err := s.api.AnswerCallback(ctx, callbackID, text); err != nil {
		return fmt.Errorf("answer callback: %w", err)
	}
	return nil
}

func (s *Screen) DeleteIncoming(ctx context.Context, t *turn.Turn) error {
	if t.IncomingID == 0 {
		return nil
	}
	s.delete(ctx, t, t.IncomingID)
	return nil
}

func (s *Screen) DismissError(ctx context.Context, t *turn.Turn) error {
	if t.Session.ErrorMessageID == 0 {
		return nil
	}
	s.delete(ctx, t, t.Session.ErrorMessageID)
	t.Session.ErrorMessageID = 0
	return s.Save(ctx, t)
}

func (s *Screen) ClearUserNoise(ctx context.Context, t *turn.Turn) error {
	if err := s.DeleteIncoming(ctx, t); err != nil {
		return err
	}
	return s.DismissError(ctx, t)
}

func (s *Screen) Save(ctx context.Context, t *turn.Turn) error {
	if err := s.sessions.Save(ctx, t.Session); err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}

func (s *Screen) delete(ctx context.Context, t *turn.Turn, messageID int64) {
	if err := s.api.DeleteMessage(ctx, t.ChatID, messageID); err != nil {
		slog.Warn("delete message", "err", err, "telegram_id", t.Actor.TelegramID)
	}
}
