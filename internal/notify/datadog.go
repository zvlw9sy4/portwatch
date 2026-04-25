package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// DatadogNotifier sends port-change events to the Datadog Events API.
type DatadogNotifier struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewDatadogNotifier creates a notifier that posts events to Datadog.
// apiKey is the Datadog API key. baseURL defaults to the US Datadog endpoint.
func NewDatadogNotifier(apiKey, baseURL string) *DatadogNotifier {
	if baseURL == "" {
		baseURL = "https://api.datadoghq.com"
	}
	return &DatadogNotifier{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type datadogEvent struct {
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	AlertType string   `json:"alert_type"`
	Tags      []string `json:"tags"`
}

// Notify sends each alert event to Datadog as a separate event.
func (d *DatadogNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	for _, ev := range events {
		alertType := "info"
		if ev.Kind == alert.Opened {
			alertType = "warning"
		}
		payload := datadogEvent{
			Title:     fmt.Sprintf("portwatch: port %s", ev.Port),
			Text:      fmt.Sprintf("Port %s was %s", ev.Port, ev.Kind),
			AlertType: alertType,
			Tags:      []string{"source:portwatch", fmt.Sprintf("port:%s", ev.Port)},
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("datadog: marshal: %w", err)
		}
		url := fmt.Sprintf("%s/api/v1/events?api_key=%s", d.baseURL, d.apiKey)
		resp, err := d.client.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("datadog: post: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("datadog: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}
