package refresh

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/Zoomish/portfolio-backend/internal/content"
	"github.com/Zoomish/portfolio-backend/internal/content/github"
	"github.com/Zoomish/portfolio-backend/internal/content/linkedin"
)

func TestCycleSavesEachSource(t *testing.T) {
	store := &memStore{}
	cycle := NewCycle(
		fakeRepos{list: []github.Candidate{{Name: "api", FullName: "owner/api", HTMLURL: "https://github.com/owner/api"}}},
		fakeProfiles{body: map[linkedin.Part][]byte{
			linkedin.Posts:   []byte(`{"data":{"posts":[]}}`),
			linkedin.Profile: []byte(`{"data":{"basic_info":{}}}`),
			linkedin.Contact: []byte(`{"data":{"email":null}}`),
		}},
		store,
	)

	if err := cycle.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(store.saved) != 4 {
		t.Fatalf("saved %d", len(store.saved))
	}
}

func TestCycleKeepsSuccessfulSources(t *testing.T) {
	store := &memStore{}
	cycle := NewCycle(
		fakeRepos{list: []github.Candidate{{Name: "api", FullName: "owner/api"}}},
		fakeProfiles{
			body: map[linkedin.Part][]byte{
				linkedin.Posts:   []byte(`{"ok":true}`),
				linkedin.Profile: []byte(`{"ok":true}`),
			},
			err: map[linkedin.Part]error{linkedin.Contact: errors.New("timeout")},
		},
		store,
	)

	err := cycle.Run(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if len(store.saved) != 3 {
		t.Fatalf("saved %d", len(store.saved))
	}
	for _, snap := range store.saved {
		if snap.Source == content.SourceLinkedInContact {
			t.Fatal("failed source was saved")
		}
	}
}

type fakeRepos struct {
	list []github.Candidate
}

func (f fakeRepos) List(context.Context) ([]github.Candidate, error) { return f.list, nil }

func (f fakeRepos) Readme(context.Context, string) (string, error) { return "", nil }

type fakeProfiles struct {
	body map[linkedin.Part][]byte
	err  map[linkedin.Part]error
}

func (f fakeProfiles) Fetch(_ context.Context, part linkedin.Part) (json.RawMessage, error) {
	if err := f.err[part]; err != nil {
		return nil, err
	}
	return json.RawMessage(f.body[part]), nil
}

type memStore struct {
	mu    sync.Mutex
	saved []content.Snapshot
}

func (m *memStore) Save(_ context.Context, snap content.Snapshot) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.saved = append(m.saved, snap)
	return nil
}
