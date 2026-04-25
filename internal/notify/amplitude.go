package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// AmplitudeNotifier sends port-change events to Amplitude Analytics.
type AmplitudeNotifier struct {
	apiKey  string
	endpoint string
	client  *http.Client
}

type amplitudePayload struct {
	APIKey string           `json:"api_key"`
	Events []amplitudeEvent `json:"events"`
}

type amplitudeEvent struct {
	UserID      string                 `json:"user_id"`
	EventType   string                 `json:"event_type"`
	Time        int64                  `json:"time"`
	EventProps  map[string]interface{} `json:"event_properties"`
}

// NewAmplitudeNotifier creates a notifier that forwards events to Amplitude.
// endpoint defaults to the Amplitude HTTP API v2 if empty.
func NewAmplitudeNotifier(apiKey, endpoint string) *AmplitudeNotifier {
	if endpoint == "" {
		endpoint = "https://api2.amplitude.com/2/httpapi"
	}
	return &AmplitudeNotifier{
		apiKey:   apiKey,
		endpoint: endpoint,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify implements alert.Notifier.
func (a *AmplitudeNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var ampEvents []amplitudeEvent
	for _, e := range events {
		ampEvents = append(ampEvents, amplitudeEvent{
			UserID:    "portwatch",
			EventType: fmt.Sprintf("port_%s", e.Kind),
			Time:      time.Now().UnixMilli(),
			EventProps: map[string]interface{}{
				"port":     e.Port.Number,
				"protocol": e.Port.Proto,
			},
		})
	}

	body, err := json.Marshal(amplitudePayload{APIKey: a.apiKey, Events: ampEvents})
	if err != nil {
		return fmt.Errorf("amplitude: marshal: %w", err)
	}

	resp, err := a.client.Post(a.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("amplitude: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("amplitude: unexpected status %d", resp.StatusCode)
	}
	return nil
}
