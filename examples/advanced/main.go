package main

import (
	"context"
	"log"
	"time"

	"github.com/cph-dev/notif"
	"github.com/cph-dev/notif/decorator"
)

func main() {
	// Create base notifier
	slack := notif.NewSlackNotifier("https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
		// Prevents any single request from hanging too long
		notif.WithTimeout(10*time.Second),
	)

	// Wrap with decorators
	// WithLogging should be the first decorator, as it uses the Name() from prior notifier when logging
	withLogging := decorator.WithLogging(slack, nil)
	withRetry := decorator.WithRetry(withLogging, 3, time.Second, 10*time.Second)

	// Create a urgent priority message
	msg := notif.Message{
		Title:    "Critical Alert",
		Content:  "Production server is down",
		URI:      "https://status.example.com",
		Priority: notif.PriorityUrgent,
		Extra: map[string]any{
			"Server":   "prod-api-01",
			"Region":   "us-east-1",
			"Impacted": "2000 users",
		},
	}

	// Send with context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := withRetry.Send(ctx, msg); err != nil {
		log.Fatalf("Failed to send notification: %v", err)
	}

	// Example of multiple notifiers
	sendToMultiple(ctx, msg)
}

func sendToMultiple(ctx context.Context, msg notif.Message) {
	// Create multiple notifiers
	slack1 := notif.NewSlackNotifier("https://hooks.slack.com/services/TEAM1/WEBHOOK/URL")
	slack2 := notif.NewSlackNotifier("https://hooks.slack.com/services/TEAM2/WEBHOOK/URL")

	// Send to multiple channels
	notifiers := []notif.Notifier{slack1, slack2}

	for i, n := range notifiers {
		if err := n.Send(ctx, msg); err != nil {
			log.Printf("Failed to send to notifier %d: %v", i, err)
		}
	}
}
