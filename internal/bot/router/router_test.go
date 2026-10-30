package router

import (
	"context"
	"testing"

	"github.com/Zoomish/portfolio-backend/internal/bot/screen"
	"github.com/Zoomish/portfolio-backend/internal/bot/session"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
	"github.com/Zoomish/portfolio-backend/internal/domain"
	"github.com/Zoomish/portfolio-backend/internal/telegram"
)

func TestCallbackFrozenWhileWaiting(t *testing.T) {
	api := &fakeAPI{}
	sessions := &memSession{state: session.State{WaitingFor: session.WaitErrorReport}}
	callbacks := &called{}
	r := &Router{
		Users:     memUsers{},
		Sessions:  sessions,
		Screen:    screen.New(api, sessions),
		Callbacks: callbacks,
		AdminID:   9,
	}
	err := r.Handle(context.Background(), telegram.Update{
		UpdateID: 3,
		Callback: &telegram.Callback{
			ID:       "cb",
			From:     telegram.User{ID: 5, FirstName: "Ann"},
			ChatID:   5,
			ChatType: "private",
			Data:     "menu-about",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if callbacks.n != 0 {
		t.Fatal("callback navigated while waiting")
	}
	if len(api.toasts) != 1 || api.toasts[0] != session.Hint(session.WaitErrorReport) {
		t.Fatalf("toasts %v", api.toasts)
	}
}

type memUsers struct{}

func (memUsers) Upsert(_ context.Context, telegramID int64, name, username string, admin bool) (domain.User, error) {
	return domain.User{ID: 1, TelegramID: telegramID, Name: name, Username: username, Admin: admin}, nil
}

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

type called struct{ n int }

func (c *called) Handle(context.Context, *turn.Turn, string) error {
	c.n++
	return nil
}

type fakeAPI struct {
	toasts []string
}

func (f *fakeAPI) SendMessage(context.Context, int64, string, telegram.Keyboard) (int64, error) {
	return 1, nil
}

func (f *fakeAPI) EditMessageText(context.Context, int64, int64, string, telegram.Keyboard) error {
	return nil
}

func (f *fakeAPI) DeleteMessage(context.Context, int64, int64) error { return nil }

func (f *fakeAPI) AnswerCallback(_ context.Context, _, text string) error {
	f.toasts = append(f.toasts, text)
	return nil
}
