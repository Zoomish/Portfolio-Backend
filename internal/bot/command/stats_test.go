package command

import (
	"context"
	"testing"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/bot/screen"
	"github.com/Zoomish/portfolio-backend/internal/bot/session"
	"github.com/Zoomish/portfolio-backend/internal/bot/stats"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
	"github.com/Zoomish/portfolio-backend/internal/domain"
	"github.com/Zoomish/portfolio-backend/internal/telegram"
)

func TestStatsRejectsStranger(t *testing.T) {
	api := &fakeAPI{}
	sessions := &memSession{}
	h := &Handler{
		Screen:  screen.New(api, sessions),
		Stats:   failQuery{},
		AdminID: 9,
		Now:     func() time.Time { return time.Unix(0, 0) },
	}
	err := h.Handle(context.Background(), &turn.Turn{
		ChatID:  5,
		Actor:   domain.User{ID: 1, TelegramID: 5},
		Session: session.State{UserID: 1},
	}, "stats")
	if err != nil {
		t.Fatal(err)
	}
	if len(api.texts) != 1 || api.texts[0] != "Эта команда только для владельца." {
		t.Fatalf("texts %v", api.texts)
	}
}

type failQuery struct{}

func (failQuery) Summary(context.Context, time.Time) (stats.Summary, error) {
	panic("summary called")
}

func (failQuery) Recent(context.Context, int) ([]domain.User, error) {
	panic("recent called")
}

type fakeAPI struct {
	texts []string
}

func (f *fakeAPI) SendMessage(_ context.Context, _ int64, text string, _ telegram.Keyboard) (int64, error) {
	f.texts = append(f.texts, text)
	return int64(len(f.texts)), nil
}

func (f *fakeAPI) EditMessageText(context.Context, int64, int64, string, telegram.Keyboard) error {
	return nil
}

func (f *fakeAPI) DeleteMessage(context.Context, int64, int64) error { return nil }

func (f *fakeAPI) AnswerCallback(context.Context, string, string) error { return nil }

type memSession struct {
	state session.State
}

func (m *memSession) Load(_ context.Context, userID int64) (session.State, error) {
	m.state.UserID = userID
	return m.state, nil
}

func (m *memSession) Save(_ context.Context, state session.State) error {
	m.state = state
	return nil
}
