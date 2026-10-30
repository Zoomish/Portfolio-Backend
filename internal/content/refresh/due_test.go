package refresh

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/content"
)

func TestNeedsRefresh(t *testing.T) {
	now := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	fresh := content.Snapshot{Payload: []byte(`{}`), FetchedAt: now.Add(-time.Hour)}
	stale := content.Snapshot{Payload: []byte(`{}`), FetchedAt: now.Add(-7 * time.Hour)}

	tests := []struct {
		name  string
		snaps map[content.Source]content.Snapshot
		err   error
		want  bool
	}{
		{name: "missing", snaps: map[content.Source]content.Snapshot{}, want: true},
		{
			name: "fresh",
			snaps: map[content.Source]content.Snapshot{
				content.SourceGithubRepos:     fresh,
				content.SourceLinkedInPosts:   fresh,
				content.SourceLinkedInProfile: fresh,
				content.SourceLinkedInContact: fresh,
			},
		},
		{
			name: "one stale",
			snaps: map[content.Source]content.Snapshot{
				content.SourceGithubRepos:     stale,
				content.SourceLinkedInPosts:   fresh,
				content.SourceLinkedInProfile: fresh,
				content.SourceLinkedInContact: fresh,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NeedsRefresh(context.Background(), dueStore{snaps: tt.snaps, err: tt.err}, now, 6*time.Hour)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %v", got)
			}
		})
	}

	_, err := NeedsRefresh(context.Background(), dueStore{err: errors.New("db")}, now, time.Hour)
	if err == nil {
		t.Fatal("expected error")
	}
}

type dueStore struct {
	snaps map[content.Source]content.Snapshot
	err   error
}

func (s dueStore) Get(_ context.Context, source content.Source) (content.Snapshot, error) {
	if s.err != nil {
		return content.Snapshot{}, s.err
	}
	snap, ok := s.snaps[source]
	if !ok {
		return content.Snapshot{}, content.ErrNotFound
	}
	return snap, nil
}
