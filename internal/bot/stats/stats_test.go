package stats

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/domain"
)

func TestAllowed(t *testing.T) {
	if Allowed(9, 0) || Allowed(5, 9) {
		t.Fatal("stranger or unset admin")
	}
	if !Allowed(9, 9) {
		t.Fatal("owner")
	}
}

func TestLoadAndFormat(t *testing.T) {
	now := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	users := &fakeUsers{total: 4, week: 2, recent: []domain.User{{
		Name: "Egor <dev>", Username: "zoomish", TelegramID: 9,
		UpdatedAt: now,
	}}}
	events := fakeEvents{counts: map[string]int{
		domain.EventStart: 3, domain.EventAbout: 2, domain.EventProjects: 1,
	}}
	summary, err := NewQuery(users, events).Summary(context.Background(), now)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Total != 4 || summary.Week != 2 || summary.Start != 3 || summary.Contact != 0 {
		t.Fatalf("%+v", summary)
	}
	if !users.since.Equal(now.Add(-7 * 24 * time.Hour)) {
		t.Fatalf("since %s", users.since)
	}
	text := FormatSummary(summary)
	if !strings.Contains(text, "Людей: 4") || !strings.Contains(text, "Обо мне: 2") {
		t.Fatalf("summary %s", text)
	}
	recent := FormatRecent(users.recent)
	if !strings.Contains(recent, "Egor &lt;dev&gt;") || !strings.Contains(recent, "@zoomish") || !strings.Contains(recent, "<code>9</code>") {
		t.Fatalf("recent %s", recent)
	}
}

type fakeUsers struct {
	total  int
	week   int
	since  time.Time
	recent []domain.User
}

func (f *fakeUsers) Count(context.Context) (int, error) { return f.total, nil }

func (f *fakeUsers) CountSince(_ context.Context, since time.Time) (int, error) {
	f.since = since
	return f.week, nil
}

func (f *fakeUsers) Recent(context.Context, int) ([]domain.User, error) { return f.recent, nil }

type fakeEvents struct {
	counts map[string]int
}

func (f fakeEvents) Counts(context.Context) (map[string]int, error) { return f.counts, nil }
