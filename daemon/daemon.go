package daemon

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gitgogit/config"
	"gitgogit/mirror"
)

// withRetry calls fn up to attempts times with exponential backoff.
// Backoff sequence: base, 2*base, 4*base, … capped at 5 minutes.
// Returns nil on first success, or the last error wrapped with attempt count.
func withRetry(ctx context.Context, attempts int, base time.Duration, fn func() error) error {
	const maxBackoff = 5 * time.Minute
	var err error
	for i := 0; i < attempts; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = fn()
		if err == nil {
			return nil
		}

		if i < attempts-1 {
			backoff := base * (1 << i)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}
	}
	return fmt.Errorf("failed after %d attempt(s): %w", attempts, err)
}

// SyncRepo mirrors one repo with retry. It constructs a Runner, calls Sync
// inside withRetry, and returns the results from the final attempt.
func SyncRepo(ctx context.Context, repo config.RepoConfig, logger *slog.Logger, attempts int, backoff time.Duration) []mirror.SyncResult {
	runner := mirror.NewRunner(repo, logger)
	var results []mirror.SyncResult

	withRetry(ctx, attempts, backoff, func() error {
		results = runner.Sync(ctx)
		for _, r := range results {
			if r.Err != nil {
				return r.Err
			}
		}
		return nil
	})

	return results
}
