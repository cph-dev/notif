package notif

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackNotifier sends notifications to Slack via webhook
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
	config     Config
}

// SlackMessage represents the Slack webhook payload
type slackMessage struct {
	Attachments []slackAttachment `json:"attachments,omitempty"`
}

type slackAttachment struct {
	Color     string       `json:"color,omitempty"`
	Title     string       `json:"title,omitempty"`
	TitleLink string       `json:"title_link,omitempty"`
	Text      string       `json:"text,omitempty"`
	Footer    string       `json:"footer,omitempty"`
	Timestamp int64        `json:"ts,omitempty"`
	Fields    []slackField `json:"fields,omitempty"`
}

type slackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// NewSlackNotifier creates a new Slack notifier
func NewSlackNotifier(webhookURL string, opts ...Option) *SlackNotifier {
	config := Config{
		Timeout: 10 * time.Second,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return &SlackNotifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
	}
}

// Send sends a notification to Slack
func (s *SlackNotifier) Send(ctx context.Context, msg Message) error {
	slackMsg := s.buildSlackMessage(msg)

	payload, err := json.Marshal(slackMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *SlackNotifier) Name() string {
	return "Slack"
}

func (s *SlackNotifier) buildSlackMessage(msg Message) slackMessage {
	color := s.getColorForPriority(msg.Priority)

	attachment := slackAttachment{
		Color:     color,
		Title:     msg.Title,
		TitleLink: msg.URI,
		Text:      msg.Content,
		Timestamp: time.Now().Unix(),
		Fields:    []slackField{},
	}

	// Add extra fields if present
	for k, v := range msg.Extra {
		attachment.Fields = append(attachment.Fields, slackField{
			Title: k,
			Value: fmt.Sprintf("%v", v),
			Short: true,
		})
	}

	return slackMessage{
		Attachments: []slackAttachment{attachment},
	}
}

func (s *SlackNotifier) getColorForPriority(priority Priority) string {
	switch priority {
	case PriorityLow:
		return "#3AA3E3" // Blue
	case PriorityNormal:
		return "#36a64f" // Green
	case PriorityHigh:
		return "#ff9900" // Orange
	case PriorityUrgent:
		return "#ff0000" // Red
	default:
		return "#3AA3E3" // Blue
	}
}
