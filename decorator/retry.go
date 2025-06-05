package decorator

import (
	"context"
	"fmt"
	"time"

	"github.com/cph-dev/notif"
)

// RetryNotifier adds retry logic to any notifier
type RetryNotifier struct {
	notifier      notif.Notifier
	maxRetries    int
	retryDelay    time.Duration
	maxRetryDelay time.Duration
}

// WithRetry wraps a notifier with retry logic
func WithRetry(notifier notif.Notifier, maxRetries int, retryDelay, maxRetryDelay time.Duration) *RetryNotifier {
	return &RetryNotifier{
		notifier:      notifier,
		maxRetries:    maxRetries,
		retryDelay:    retryDelay,
		maxRetryDelay: maxRetryDelay,
	}
}

// Send sends a notification with retry logic
func (r *RetryNotifier) Send(ctx context.Context, msg notif.Message) error {
	var lastErr error
	delay := r.retryDelay

	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}

			// Exponential backoff
			delay *= 2
			if delay > r.maxRetryDelay {
				delay = r.maxRetryDelay
			}
		}

		err := r.notifier.Send(ctx, msg)
		if err == nil {
			return nil
		}

		lastErr = err
	}

	return fmt.Errorf("failed after %d retries: %w", r.maxRetries, lastErr)
}

func (r *RetryNotifier) Name() string {
	return "Retry"
}
