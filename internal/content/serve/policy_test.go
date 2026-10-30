package serve

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/content"
)

func TestPolicyFreshSnapshotDoesNotRefresh(t *testing.T) {
	now := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	store := &fakeStore{snaps: map[content.Source]content.Snapshot{
		content.SourceGithubRepos: {Source: content.SourceGithubRepos, Payload: []byte(`[]`), FetchedAt: now.Add(-time.Hour)},
	}}
	refresh := &fakeRefresh{}
	svc := New(store, refresh, fakeClock{now}, 6*time.Hour)

	body, err := svc.Github(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !body.Cached || string(body.Data) != `[]` {
		t.Fatalf("body %+v", body)
	}
	if refresh.runs != 0 || refresh.triggers != 0 {
		t.Fatalf("runs %d triggers %d", refresh.runs, refresh.triggers)
	}
}

func TestPolicyStaleSnapshotTriggersBackgroundRefresh(t *testing.T) {
	now := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	store := &fakeStore{snaps: map[content.Source]content.Snapshot{
		content.SourceGithubRepos: {Source: content.SourceGithubRepos, Payload: []byte(`[{"name":"old"}]`), FetchedAt: now.Add(-7 * time.Hour)},
	}}
	refresh := &fakeRefresh{}
	svc := New(store, refresh, fakeClock{now}, 6*time.Hour)

	body, err := svc.Github(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if string(body.Data) != `[{"name":"old"}]` {
		t.Fatalf("data %s", body.Data)
	}
	if refresh.runs != 0 || refresh.triggers != 1 {
		t.Fatalf("runs %d triggers %d", refresh.runs, refresh.triggers)
	}
}

func TestPolicyMissingSnapshotRefreshesOnce(t *testing.T) {
	store := &fakeStore{snaps: map[content.Source]content.Snapshot{}}
	refresh := &fakeRefresh{onRun: func() {
		raw := []byte(`{"success":true}`)
		now := time.Now().UTC()
		store.snaps[content.SourceLinkedInPosts] = content.Snapshot{Payload: raw, FetchedAt: now}
		store.snaps[content.SourceLinkedInProfile] = content.Snapshot{Payload: raw, FetchedAt: now}
		store.snaps[content.SourceLinkedInContact] = content.Snapshot{Payload: raw, FetchedAt: now}
	}}
	svc := New(store, refresh, fakeClock{time.Now()}, time.Hour)

	body, err := svc.LinkedInAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if refresh.runs != 1 {
		t.Fatalf("runs %d", refresh.runs)
	}
	if !body.Posts.Cached || string(body.Details.Data) != `{"success":true}` {
		t.Fatalf("body %+v", body)
	}
}

func TestPolicyMissingAfterFailedRefresh(t *testing.T) {
	store := &fakeStore{snaps: map[content.Source]content.Snapshot{}}
	refresh := &fakeRefresh{err: errors.New("upstream")}
	svc := New(store, refresh, SystemClock{}, time.Hour)

	_, err := svc.Github(context.Background())
	if !errors.Is(err, content.ErrNotReady) {
		t.Fatalf("got %v", err)
	}
}

type fakeClock struct{ t time.Time }

func (f fakeClock) Now() time.Time { return f.t }

type fakeStore struct {
	snaps map[content.Source]content.Snapshot
}

func (f *fakeStore) Get(_ context.Context, source content.Source) (content.Snapshot, error) {
	snap, ok := f.snaps[source]
	if !ok {
		return content.Snapshot{}, content.ErrNotFound
	}
	return snap, nil
}

type fakeRefresh struct {
	runs     int
	triggers int
	err      error
	onRun    func()
}

func (f *fakeRefresh) RunOnce(context.Context) error {
	f.runs++
	if f.onRun != nil {
		f.onRun()
	}
	return f.err
}

func (f *fakeRefresh) Trigger() { f.triggers++ }
