package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// DiscordNotifier sends alerts to a Discord webhook.
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewDiscordNotifier creates a notifier that posts to the given Discord webhook URL.
func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

type discordPayload struct {
	Content string `json:"content"`
}

// Notify sends the event to Discord as a plain message.
func (d *DiscordNotifier) Notify(e alert.Event) error {
	var msg string
	switch e.Kind {
	case alert.Opened:
		msg = fmt.Sprintf(":warning: Port opened: %s", e.Port)
	case alert.Closed:
		msg = fmt.Sprintf(":information_source: Port closed: %s", e.Port)
	default:
		msg = fmt.Sprintf("Port event [%s]: %s", e.Kind, e.Port)
	}

	body, err := json.Marshal(discordPayload{Content: msg})
	if err != nil {
		return fmt.Errorf("discord: marshal: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}
	return nil
}
