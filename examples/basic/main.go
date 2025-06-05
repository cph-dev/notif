package main

import (
	"context"
	"log"
	"time"

	"github.com/cph-dev/notif"
)

func main() {
	// Create a Slack notifier
	slack := notif.NewSlackNotifier("https://hooks.slack.com/services/YOUR/WEBHOOK/URL")

	// Create a message
	msg := notif.Message{
		Title:    "System Alert",
		Content:  "The database backup completed successfully",
		URI:      "https://example.com/backups/latest",
		Priority: notif.PriorityNormal,
		Extra: map[string]any{
			"Duration": "2m 34s",
			"Size":     "1.2GB",
		},
	}

	// Send the notification
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := slack.Send(ctx, msg); err != nil {
		log.Fatalf("Failed to send notification: %v", err)
	}

	log.Println("Notification sent successfully")
}
