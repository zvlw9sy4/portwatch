package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// TeamsNotifier sends alerts to a Microsoft Teams channel via an incoming webhook.
type TeamsNotifier struct {
	webhookURL string
	client     *http.Client
}

type teamsPayload struct {
	Text string `json:"text"`
}

// NewTeamsNotifier returns a TeamsNotifier that posts to the given webhook URL.
func NewTeamsNotifier(webhookURL string) *TeamsNotifier {
	return &TeamsNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Notify sends a single event to the Teams channel.
func (t *TeamsNotifier) Notify(e alert.Event) error {
	body, err := json.Marshal(teamsPayload{
		Text: fmt.Sprintf("**portwatch** — %s port %s (%s)", e.Kind, e.Port.String(), e.Port.Proto),
	})
	if err != nil {
		return fmt.Errorf("teams: marshal: %w", err)
	}

	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}
