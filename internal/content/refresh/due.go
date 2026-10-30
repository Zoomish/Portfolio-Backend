package refresh

import (
	"context"
	"errors"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/content"
)

type SnapshotReader interface {
	Get(ctx context.Context, source content.Source) (content.Snapshot, error)
}

func NeedsRefresh(ctx context.Context, snaps SnapshotReader, now time.Time, maxAge time.Duration) (bool, error) {
	for _, source := range []content.Source{
		content.SourceGithubRepos,
		content.SourceLinkedInPosts,
		content.SourceLinkedInProfile,
		content.SourceLinkedInContact,
	} {
		snap, err := snaps.Get(ctx, source)
		if errors.Is(err, content.ErrNotFound) {
			return true, nil
		}
		if err != nil {
			return false, err
		}
		if content.IsStale(now, snap.FetchedAt, maxAge) {
			return true, nil
		}
	}
	return false, nil
}
