package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// GoogleChatNotifier sends port change alerts to a Google Chat webhook.
type GoogleChatNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewGoogleChatNotifier creates a notifier that posts messages to the given
// Google Chat incoming webhook URL.
func NewGoogleChatNotifier(webhookURL string) *GoogleChatNotifier {
	return &GoogleChatNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

type googleChatPayload struct {
	Text string `json:"text"`
}

// Notify sends all events to Google Chat as a single formatted message.
func (g *GoogleChatNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		switch e.Kind {
		case alert.Opened:
			fmt.Fprintf(&buf, "🟢 Port opened: %s\n", e.Port)
		case alert.Closed:
			fmt.Fprintf(&buf, "🔴 Port closed: %s\n", e.Port)
		}
	}

	payload := googleChatPayload{Text: buf.String()}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: marshal payload: %w", err)
	}

	resp, err := g.client.Post(g.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
