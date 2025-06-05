package decorator

import (
	"context"
	"log"
	"time"

	"github.com/cph-dev/notif"
)

// LoggingNotifier adds logging to any notifier
type LoggingNotifier struct {
	notifier notif.Notifier
	logger   *log.Logger
}

// WithLogging wraps a notifier with logging.
//
// Should be placed right after base notifier as the Name() of prior notifier is used when logging.
func WithLogging(notifier notif.Notifier, logger *log.Logger) *LoggingNotifier {
	if logger == nil {
		logger = log.Default()
	}
	return &LoggingNotifier{
		notifier: notifier,
		logger:   logger,
	}
}

// Send sends a notification with logging
func (l *LoggingNotifier) Send(ctx context.Context, msg notif.Message) error {
	start := time.Now()

	l.logger.Printf("Sending notification via %s: Title=%s, Priority=%d", l.notifier.Name(), msg.Title, msg.Priority)

	err := l.notifier.Send(ctx, msg)

	duration := time.Since(start)

	if err != nil {
		l.logger.Printf("Notification failed: %v (duration: %s)", err, duration)
	} else {
		l.logger.Printf("Notification sent successfully (duration: %s)", duration)
	}

	return err
}

func (l *LoggingNotifier) Name() string {
	return "Logging"
}
