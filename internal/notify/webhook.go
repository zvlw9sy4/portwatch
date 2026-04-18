package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier sends alert events to an HTTP endpoint as JSON.
type WebhookNotifier struct {
	URL    string
	client *http.Client
}

// WebhookPayload is the JSON body sent on each alert.
type WebhookPayload struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
}

// NewWebhookNotifier creates a WebhookNotifier with a default timeout.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL: url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Notify sends a single payload to the configured webhook URL.
func (w *WebhookNotifier) Notify(payload WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook marshal: %w", err)
	}
	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
