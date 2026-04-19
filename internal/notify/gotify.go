package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// GotifyNotifier sends alerts to a Gotify server.
type GotifyNotifier struct {
	baseURL  string
	token    string
	client   *http.Client
	priority int
}

type gotifyPayload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

// NewGotifyNotifier creates a notifier that pushes messages to Gotify.
// baseURL should be the root URL of the Gotify server (e.g. "http://localhost:80").
func NewGotifyNotifier(baseURL, token string, priority int, client *http.Client) *GotifyNotifier {
	if client == nil {
		client = http.DefaultClient
	}
	return &GotifyNotifier{baseURL: baseURL, token: token, priority: priority, client: client}
}

// Notify sends a single alert event to the Gotify server.
func (g *GotifyNotifier) Notify(e alert.Event) error {
	title := fmt.Sprintf("Port %s", e.Kind)
	msg := fmt.Sprintf("Port %d/%s %s", e.Port.Number, e.Port.Protocol, e.Kind)

	p := gotifyPayload{Title: title, Message: msg, Priority: g.priority}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("gotify: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/message?token=%s", g.baseURL, g.token)
	resp, err := g.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gotify: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
