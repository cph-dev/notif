package notif

import (
	"net/http"
	"time"
)

// WithTimeout sets the timeout for notifications
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// SlackOption is a function that modifies SlackNotifier configuration
type SlackOption func(*SlackNotifier)

// WithHTTPClient sets a custom HTTP client for Slack notifications
func WithHTTPClient(client *http.Client) SlackOption {
	return func(s *SlackNotifier) {
		s.client = client
	}
}
