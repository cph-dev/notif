package notif

import (
	"context"
	"time"
)

// Message represents a normalized notification message.
// All notifiers expect messages in this format.
type Message struct {
	Title    string         // Subject line or header
	Content  string         // Main content body
	URI      string         // Optional link for more details
	Priority Priority       // Message priority level
	Extra    map[string]any // Additional fields for rich notifications
}

// Priority represents the urgency level of a notification
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
	PriorityUrgent
)

// Notifier is the interface that all notification methods must implement
type Notifier interface {
	Send(ctx context.Context, msg Message) error
	Name() string
}

// Config holds common configuration options
type Config struct {
	Timeout time.Duration
}

// Option is a function that modifies configuration
type Option func(*Config)
