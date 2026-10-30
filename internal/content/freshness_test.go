package content

import (
	"testing"
	"time"
)

func TestIsStale(t *testing.T) {
	now := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	ttl := 6 * time.Hour

	tests := []struct {
		name    string
		fetched time.Time
		want    bool
	}{
		{name: "fresh", fetched: now.Add(-time.Hour), want: false},
		{name: "exact ttl", fetched: now.Add(-ttl), want: true},
		{name: "older", fetched: now.Add(-7 * time.Hour), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStale(now, tt.fetched, ttl); got != tt.want {
				t.Fatalf("got %v", got)
			}
		})
	}
}
