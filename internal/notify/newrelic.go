package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

const defaultNewRelicURL = "https://log-api.newrelic.com/log/v1"

// NewRelicNotifier sends port change events to New Relic Logs API.
type NewRelicNotifier struct {
	apiKey  string
	endpoint string
	client  *http.Client
}

// NewNewRelicNotifier creates a notifier that forwards events to New Relic.
// apiKey is the New Relic Ingest License Key. endpoint may be empty to use
// the default US Log API URL.
func NewNewRelicNotifier(apiKey, endpoint string) *NewRelicNotifier {
	if endpoint == "" {
		endpoint = defaultNewRelicURL
	}
	return &NewRelicNotifier{
		apiKey:   apiKey,
		endpoint: endpoint,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

type newRelicLog struct {
	Timestamp int64             `json:"timestamp"`
	Message   string            `json:"message"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type newRelicPayload struct {
	Logs []newRelicLog `json:"logs"`
}

// Notify sends all events to New Relic as structured log entries.
func (n *NewRelicNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	logs := make([]newRelicLog, 0, len(events))
	for _, e := range events {
		logs = append(logs, newRelicLog{
			Timestamp: time.Now().UnixMilli(),
			Message:   fmt.Sprintf("portwatch: port %s %s", e.Port, e.Kind),
			Attributes: map[string]string{
				"port":     fmt.Sprintf("%d", e.Port.Number),
				"protocol": e.Port.Proto,
				"kind":     string(e.Kind),
				"source":   "portwatch",
			},
		})
	}

	body, err := json.Marshal(newRelicPayload{Logs: logs})
	if err != nil {
		return fmt.Errorf("newrelic: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("newrelic: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", n.apiKey)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("newrelic: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("newrelic: unexpected status %d", resp.StatusCode)
	}
	return nil
}
