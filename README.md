# notif

A simple, extensible notification library for GO applications using only the standard library.

## Installation

```bash
go get github.com/cph-dev/notif
```

## Basic Usage

```go
package main

import (
    "context"
    "github.com/cph-dev/notif"
)

func main() {
    // Create a Slack notifier
    slack := notif.NewSlackNotifier("https://hooks.slack.com/services/YOUR/WEBHOOK/URL")

    // Create a message
    msg := notif.Message{
        Title:       "System Alert",
        Description: "Database backup completed",
        Priority:    notif.PriorityNormal,
    }

    // Send notification
    if err := slack.Send(context.Background(), msg); err != nil {
        // Handle error
    }
}
```

## Advanced Usage

### With Decorators

```go
import "github.com/cph-dev/notif/decorator"

// Add logging
withLogging := decorator.WithLogging(withRetry, nil)

// Add retry logic
withRetry := decorator.WithRetry(slack, 3, time.Second, 10*time.Second)


// Use the decorated notifier
err := withRetry.Send(ctx, msg)
```

### Multiple Notifiers

```go
// Create multiple notifiers
slack := notif.NewSlackNotifier(slackURL)
email := notif.NewEmailNotifier(emailConfig) // Future

// Send to both
for _, n := range []notif.Notifier{slack, email} {
	if err := n.Send(ctx, msg); err != nil {
		// Handle error
	}
}
```

## Message Format

```go
type Message struct {
    Title       string              // Subject or header
    Content     string              // Main content
    URI         string              // Optional link
    Priority    Priority            // Low, Normal, High, Urgent
    Extra       map[string]any      // Additional fields
}
```

## Extending

Implement the `Notifier` interface:

```go
type Notifier interface {
    Send(ctx context.Context, msg Message) error
    Name() string
}
```

## Features

- **Composable**: Combine with decorators for additional functionality
- **Zero Dependencies**: Uses only Go standard library
- **Context Support**: Full support for cancellation and timeouts

## License

MIT
