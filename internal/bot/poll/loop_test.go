package poll

import (
	"context"
	"testing"

	"github.com/Zoomish/portfolio-backend/internal/telegram"
)

func TestLoopAdvancesOffset(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var offsets []int64
	var seen int64
	loop := New(func(_ context.Context, offset int64) ([]telegram.Update, error) {
		offsets = append(offsets, offset)
		if offset == 0 {
			return []telegram.Update{{UpdateID: 4}}, nil
		}
		cancel()
		return nil, context.Canceled
	}, func(_ context.Context, update telegram.Update) error {
		seen = update.UpdateID
		return nil
	}, nil, nil, 0)
	loop.Start(ctx)
	if seen != 4 {
		t.Fatalf("seen %d", seen)
	}
	if len(offsets) != 2 || offsets[0] != 0 || offsets[1] != 5 {
		t.Fatalf("offsets %v", offsets)
	}
}
