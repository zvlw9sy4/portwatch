package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// MattermostNotifier sends port change alerts to a Mattermost incoming webhook.
type MattermostNotifier struct {
	webhookURL string
	channel    string
	client     *http.Client
}

type mattermostPayload struct {
	Channel string `json:"channel,omitempty"`
	Text    string `json:"text"`
}

// NewMattermostNotifier creates a notifier that posts to the given Mattermost
// incoming webhook URL. channel may be empty to use the webhook's default.
func NewMattermostNotifier(webhookURL, channel string) *MattermostNotifier {
	return &MattermostNotifier{
		webhookURL: webhookURL,
		channel:    channel,
		client:     &http.Client{},
	}
}

// Notify sends all events as a single Mattermost message.
func (m *MattermostNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		switch e.Kind {
		case alert.Opened:
			fmt.Fprintf(&buf, ":large_green_circle: Port **%s** opened\n", e.Port)
		case alert.Closed:
			fmt.Fprintf(&buf, ":red_circle: Port **%s** closed\n", e.Port)
		}
	}

	payload := mattermostPayload{
		Channel: m.channel,
		Text:    buf.String(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: marshal payload: %w", err)
	}

	resp, err := m.client.Post(m.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mattermost: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status %d", resp.StatusCode)
	}
	return nil
}
