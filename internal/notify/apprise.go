package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// AppriseNotifier sends alerts to an Apprise API server.
// Apprise supports 80+ notification services through a unified API.
// See https://github.com/caronc/apprise-api for server setup.
type AppriseNotifier struct {
	url   string
	tag   string
	client *http.Client
}

type apprisePayload struct {
	URLs    string `json:"urls,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Type    string `json:"type"`
}

// NewAppriseNotifier creates an AppriseNotifier that posts to the given
// Apprise API base URL. The tag parameter optionally scopes delivery to
// a pre-configured tag group on the server; leave empty to use the
// server default.
func NewAppriseNotifier(baseURL, tag string) *AppriseNotifier {
	return &AppriseNotifier{
		url:    baseURL + "/notify",
		tag:    tag,
		client: &http.Client{},
	}
}

// Notify implements alert.Notifier.
func (n *AppriseNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	title, body := buildSummary(events)

	p := apprisePayload{
		Title: title,
		Body:  body,
		Type:  "warning",
	}
	if n.tag != "" {
		p.Tag = n.tag
	}

	raw, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("apprise: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.url, "application/json", bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("apprise: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("apprise: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// buildSummary returns a short title and a multi-line body for the event list.
func buildSummary(events []alert.Event) (string, string) {
	title := fmt.Sprintf("portwatch: %d port change(s) detected", len(events))
	var buf bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&buf, "[%s] %s\n", e.Kind, e.Port)
	}
	return title, buf.String()
}
