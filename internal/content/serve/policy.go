package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/content"
)

const syncTimeout = 45 * time.Second

type SnapshotStore interface {
	Get(ctx context.Context, source content.Source) (content.Snapshot, error)
}

type RefreshControl interface {
	RunOnce(ctx context.Context) error
	Trigger()
}

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now() }

type Service struct {
	snaps   SnapshotStore
	refresh RefreshControl
	clock   Clock
	maxAge  time.Duration
}

func New(snaps SnapshotStore, refresh RefreshControl, clock Clock, maxAge time.Duration) *Service {
	return &Service{snaps: snaps, refresh: refresh, clock: clock, maxAge: maxAge}
}

func (s *Service) ensure(ctx context.Context, sources ...content.Source) error {
	missing := false
	stale := false
	for _, source := range sources {
		snap, err := s.snaps.Get(ctx, source)
		if err != nil {
			if errors.Is(err, content.ErrNotFound) {
				missing = true
				continue
			}
			return fmt.Errorf("ensure %s: %w", source, err)
		}
		if content.IsStale(s.clock.Now(), snap.FetchedAt, s.maxAge) {
			stale = true
		}
	}
	if missing {
		syncCtx, cancel := context.WithTimeout(ctx, syncTimeout)
		defer cancel()
		if err := s.refresh.RunOnce(syncCtx); err != nil {
			slog.Error("sync refresh", "err", err)
		}
		return nil
	}
	if stale {
		s.refresh.Trigger()
	}
	return nil
}

func (s *Service) read(ctx context.Context, source content.Source) (content.Snapshot, error) {
	snap, err := s.snaps.Get(ctx, source)
	if errors.Is(err, content.ErrNotFound) {
		return content.Snapshot{}, content.ErrNotReady
	}
	if err != nil {
		return content.Snapshot{}, fmt.Errorf("read %s: %w", source, err)
	}
	if !json.Valid(snap.Payload) {
		return content.Snapshot{}, fmt.Errorf("%s snapshot: invalid json", source)
	}
	return snap, nil
}
