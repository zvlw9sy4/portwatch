package notify

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// NtfyNotifier sends alerts to an ntfy.sh topic via HTTP.
// See https://ntfy.sh for details.
type NtfyNotifier struct {
	serverURL string
	topic     string
	client    *http.Client
}

// NewNtfyNotifier creates a new NtfyNotifier.
// serverURL should be the base URL (e.g. "https://ntfy.sh" or a self-hosted instance).
// topic is the ntfy topic to publish to.
func NewNtfyNotifier(serverURL, topic string) *NtfyNotifier {
	return &NtfyNotifier{
		serverURL: strings.TrimRight(serverURL, "/"),
		topic:     topic,
		client:    &http.Client{},
	}
}

// Notify sends a single alert event to the configured ntfy topic.
func (n *NtfyNotifier) Notify(event alert.Event) error {
	var action string
	switch event.Kind {
	case alert.Opened:
		action = "opened"
	case alert.Closed:
		action = "closed"
	default:
		action = "changed"
	}

	title := fmt.Sprintf("Port %s", action)
	body := fmt.Sprintf("Port %d/%s is now %s",
		event.Port.Number, event.Port.Protocol, action)

	url := fmt.Sprintf("%s/%s", n.serverURL, n.topic)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("ntfy: build request: %w", err)
	}
	req.Header.Set("Title", title)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("ntfy: send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status %d", resp.StatusCode)
	}
	return nil
}
