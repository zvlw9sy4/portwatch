package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// SlackNotifier sends alert events to a Slack webhook URL.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewSlackNotifier creates a SlackNotifier that posts to the given Slack webhook URL.
func NewSlackNotifier(webhookURL string, client *http.Client) *SlackNotifier {
	if client == nil {
		client = http.DefaultClient
	}
	return &SlackNotifier{webhookURL: webhookURL, client: client}
}

type slackPayload struct {
	Text string `json:"text"`
}

// Notify sends a Slack message for the given event.
func (s *SlackNotifier) Notify(e alert.Event) error {
	text := fmt.Sprintf("[portwatch] %s port %s/%s",
		e.Kind, e.Port.Address, e.Port.Proto)

	body, err := json.Marshal(slackPayload{Text: text})
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
