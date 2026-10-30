package refresh

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunOnceJoinsInFlightRefresh(t *testing.T) {
	var calls atomic.Int32
	started := make(chan struct{})
	release := make(chan struct{})
	loop := NewLoop(func(context.Context) error {
		if calls.Add(1) == 1 {
			close(started)
		}
		<-release
		return nil
	}, time.Hour, 3, nil, nil)

	var wg sync.WaitGroup
	wg.Add(2)
	errCh := make(chan error, 2)
	go func() {
		defer wg.Done()
		errCh <- loop.RunOnce(context.Background())
	}()
	<-started
	go func() {
		defer wg.Done()
		errCh <- loop.RunOnce(context.Background())
	}()
	time.Sleep(30 * time.Millisecond)
	close(release)
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("refresh calls %d", got)
	}
}

func TestStartSkipsRefreshWhenFresh(t *testing.T) {
	var calls atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	loop := NewLoop(func(context.Context) error {
		calls.Add(1)
		return nil
	}, time.Hour, 3, nil, func(context.Context) bool { return false })

	done := make(chan struct{})
	go func() {
		loop.Start(ctx)
		close(done)
	}()
	time.Sleep(30 * time.Millisecond)
	cancel()
	<-done
	if calls.Load() != 0 {
		t.Fatalf("refresh calls %d", calls.Load())
	}
}

func TestTriggerStopsWhenProcessContextEnds(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{})
	finished := make(chan struct{})
	loop := NewLoop(func(ctx context.Context) error {
		close(started)
		<-ctx.Done()
		close(finished)
		return ctx.Err()
	}, time.Hour, 3, nil, nil)
	loop.Bind(ctx)
	loop.Trigger()
	<-started
	cancel()
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("refresh kept running after cancel")
	}
}
