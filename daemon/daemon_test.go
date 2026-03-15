package daemon

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTest = errors.New("test error")

func TestWithRetry_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	err := withRetry(context.Background(), 3, time.Millisecond, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestWithRetry_SucceedsOnSecondAttempt(t *testing.T) {
	calls := 0
	err := withRetry(context.Background(), 3, time.Millisecond, func() error {
		calls++
		if calls < 2 {
			return errTest
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected nil after second attempt, got %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestWithRetry_ExhaustsAllAttempts(t *testing.T) {
	calls := 0
	err := withRetry(context.Background(), 3, time.Millisecond, func() error {
		calls++
		return errTest
	})
	if err == nil {
		t.Error("expected error after all attempts")
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
	if !errors.Is(err, errTest) {
		t.Errorf("expected wrapped errTest, got %v", err)
	}
}

func TestWithRetry_ExponentialBackoff(t *testing.T) {
	base := 10 * time.Millisecond
	calls := 0
	timestamps := []time.Time{}

	withRetry(context.Background(), 4, base, func() error {
		timestamps = append(timestamps, time.Now())
		calls++
		return errTest
	})

	if calls != 4 {
		t.Fatalf("expected 4 calls, got %d", calls)
	}

	// Gap between attempt 1→2 should be ~base (10ms)
	// Gap between attempt 2→3 should be ~2*base (20ms)
	// Gap between attempt 3→4 should be ~4*base (40ms)
	gaps := []time.Duration{
		timestamps[1].Sub(timestamps[0]),
		timestamps[2].Sub(timestamps[1]),
		timestamps[3].Sub(timestamps[2]),
	}
	expected := []time.Duration{base, 2 * base, 4 * base}

	for i, gap := range gaps {
		// Allow 5x margin for slow CI environments.
		if gap < expected[i] || gap > expected[i]*5 {
			t.Errorf("gap[%d] = %v, want ~%v", i, gap, expected[i])
		}
	}
}

func TestWithRetry_ContextCancelledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := withRetry(ctx, 3, time.Millisecond, func() error {
		calls++
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Errorf("expected 0 calls, got %d", calls)
	}
}

func TestWithRetry_ContextCancelledDuringSleep(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	// Use a long backoff so the cancel fires during the sleep.
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	err := withRetry(ctx, 5, time.Second, func() error {
		calls++
		return errTest
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	// Should have made exactly 1 call before being cancelled during sleep.
	if calls != 1 {
		t.Errorf("expected 1 call before cancel, got %d", calls)
	}
}
