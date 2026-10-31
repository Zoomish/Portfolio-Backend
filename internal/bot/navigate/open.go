package navigate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Zoomish/portfolio-backend/internal/bot/boterr"
	"github.com/Zoomish/portfolio-backend/internal/bot/card"
	"github.com/Zoomish/portfolio-backend/internal/bot/screen"
	"github.com/Zoomish/portfolio-backend/internal/bot/turn"
	"github.com/Zoomish/portfolio-backend/internal/content"
)

const (
	ModeCommand = "command"
	ModeEdit    = "edit"
)

type Snapshots interface {
	Profile(ctx context.Context) (json.RawMessage, error)
	Contact(ctx context.Context) (json.RawMessage, error)
	Projects(ctx context.Context) (json.RawMessage, error)
	Posts(ctx context.Context) (json.RawMessage, error)
}

type Events interface {
	Add(ctx context.Context, userID int64, kind string) error
}

type Nav struct {
	Screen *screen.Screen
	Snaps  Snapshots
	Events Events
	Alert  func(context.Context, string)
}

func (n *Nav) Open(ctx context.Context, t *turn.Turn, mode, kind string) error {
	view, err := n.render(ctx, kind)
	if errors.Is(err, content.ErrNotReady) || errors.Is(err, card.ErrIncomplete) {
		view = card.Unavailable()
	} else if err != nil {
		return boterr.Report(ctx, n.Alert, n.Screen, t, err)
	}
	if err := n.present(ctx, t, mode, view); err != nil {
		return err
	}
	if view.Event == "" {
		return nil
	}
	if err := n.Events.Add(ctx, t.Actor.ID, view.Event); err != nil {
		slog.Error("record event", "err", err, "telegram_id", t.Actor.TelegramID, "update_id", t.UpdateID)
		if n.Alert != nil {
			n.Alert(ctx, fmt.Sprintf("bot event: %v", err))
		}
	}
	return nil
}

func (n *Nav) render(ctx context.Context, kind string) (card.View, error) {
	switch kind {
	case card.ScreenHome:
		raw, err := n.Snaps.Profile(ctx)
		if err != nil {
			return card.View{}, err
		}
		return card.Home(raw)
	case card.ScreenAbout:
		raw, err := n.Snaps.Profile(ctx)
		if err != nil {
			return card.View{}, err
		}
		return card.About(raw)
	case card.ScreenProjects:
		raw, err := n.Snaps.Projects(ctx)
		if err != nil {
			return card.View{}, err
		}
		return card.Projects(raw)
	case card.ScreenPosts:
		raw, err := n.Snaps.Posts(ctx)
		if err != nil {
			return card.View{}, err
		}
		return card.Posts(raw)
	case card.ScreenContact:
		raw, err := n.Snaps.Contact(ctx)
		if err != nil {
			return card.View{}, err
		}
		return card.Contact(raw)
	default:
		return card.View{}, fmt.Errorf("unknown screen %s", kind)
	}
}

func (n *Nav) present(ctx context.Context, t *turn.Turn, mode string, view card.View) error {
	var err error
	if mode == ModeCommand {
		err = n.Screen.SendCommand(ctx, t, view.Text, view.Keys)
	} else {
		err = n.Screen.EditFunctional(ctx, t, view.Text, view.Keys)
	}
	if err != nil {
		return err
	}
	return n.Screen.DismissError(ctx, t)
}
