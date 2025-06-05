package notif

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSlackNotifier_Send(t *testing.T) {
	tests := []struct {
		name           string
		message        Message
		serverResponse int
		serverDelay    time.Duration
		wantErr        bool
	}{
		{
			name: "successful send",
			message: Message{
				Title:    "Test Alert",
				Content:  "Test description",
				Priority: PriorityNormal,
			},
			serverResponse: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "server error",
			message: Message{
				Title:   "Test Alert",
				Content: "Test description",
			},
			serverResponse: http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "timeout",
			message: Message{
				Title: "Test Alert",
			},
			serverResponse: http.StatusOK,
			serverDelay:    2 * time.Second,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedPayload slackMessage
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Simulate slow response
				if tt.serverDelay > 0 {
					time.Sleep(tt.serverDelay)
				}

				if r.Method != "POST" {
					t.Errorf("Expected POST, got %s", r.Method)
				}

				if ct := r.Header.Get("Content-Type"); ct != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", ct)
				}

				body, _ := io.ReadAll(r.Body)
				json.Unmarshal(body, &receivedPayload)

				w.WriteHeader(tt.serverResponse)
			}))
			defer server.Close()

			// Create notifier with short timeout
			notifier := NewSlackNotifier(server.URL, WithTimeout(1*time.Second))

			ctx := context.Background()
			err := notifier.Send(ctx, tt.message)

			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && err == nil {
				if receivedPayload.Attachments[0].Title != tt.message.Title {
					t.Errorf("Expected title %q, got %q", tt.message.Title, receivedPayload.Attachments[0].Title)
				}
			}
		})
	}
}

func TestSlackNotifier_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier := NewSlackNotifier(server.URL)

	// Create already cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := notifier.Send(ctx, Message{Title: "Test"})
	if err == nil {
		t.Error("Expected error with cancelled context")
	}
}

func TestSlackMessage_Building(t *testing.T) {
	notifier := &SlackNotifier{}

	msg := Message{
		Title:    "Test Title",
		Content:  "Test Description",
		URI:      "https://example.com",
		Priority: PriorityUrgent,
		Extra: map[string]any{
			"Key1": "Value1",
			"Key2": 42,
		},
	}

	slackMsg := notifier.buildSlackMessage(msg)

	if len(slackMsg.Attachments) != 1 {
		t.Fatalf("Expected 1 attachment, got %d", len(slackMsg.Attachments))
	}

	att := slackMsg.Attachments[0]
	if att.Color != "#ff0000" { // Red for urgent
		t.Errorf("Expected red color for urgent, got %s", att.Color)
	}

	if att.Title != msg.Title {
		t.Errorf("Expected %q, got %q", att.Title, msg.Title)
	}

	if len(att.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(att.Fields))
	}
}
