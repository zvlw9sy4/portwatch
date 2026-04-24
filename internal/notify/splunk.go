package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// SplunkNotifier sends port events to a Splunk HTTP Event Collector (HEC) endpoint.
type SplunkNotifier struct {
	endpoint string
	token    string
	source   string
	client   *http.Client
}

type splunkEvent struct {
	Time   float64        `json:"time"`
	Source string         `json:"source"`
	Event  map[string]any `json:"event"`
}

// NewSplunkNotifier creates a notifier that forwards events to Splunk HEC.
// endpoint should be the full HEC URL, e.g. https://splunk:8088/services/collector.
func NewSplunkNotifier(endpoint, token, source string) *SplunkNotifier {
	return &SplunkNotifier{
		endpoint: endpoint,
		token:    token,
		source:   source,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends each alert event as an individual Splunk HEC event.
func (s *SplunkNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, e := range events {
		se := splunkEvent{
			Time:   float64(e.Time.UnixNano()) / 1e9,
			Source: s.source,
			Event: map[string]any{
				"action":   e.Kind,
				"port":     e.Port.Number,
				"protocol": e.Port.Proto,
			},
		}
		if err := enc.Encode(se); err != nil {
			return fmt.Errorf("splunk: encode: %w", err)
		}
	}

	req, err := http.NewRequest(http.MethodPost, s.endpoint, &buf)
	if err != nil {
		return fmt.Errorf("splunk: build request: %w", err)
	}
	req.Header.Set("Authorization", "Splunk "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
