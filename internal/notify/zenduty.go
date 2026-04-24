package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// ZendutyNotifier sends alerts to a Zenduty integration via its Events API.
type ZendutyNotifier struct {
	integrationKey string
	client         *http.Client
}

type zendutyPayload struct {
	AlertType string                 `json:"alert_type"`
	Message   string                 `json:"message"`
	Summary   string                 `json:"summary"`
	Payload   map[string]interface{} `json:"payload"`
}

// NewZendutyNotifier creates a notifier that posts to the Zenduty Events API.
func NewZendutyNotifier(integrationKey string) *ZendutyNotifier {
	return &ZendutyNotifier{
		integrationKey: integrationKey,
		client:         &http.Client{},
	}
}

// Notify sends each event as a separate Zenduty alert.
func (z *ZendutyNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	const endpoint = "https://events.zenduty.com/api/events/"

	for _, ev := range events {
		alertType := "info"
		if ev.Type == alert.Opened {
			alertType = "critical"
		}

		p := zendutyPayload{
			AlertType: alertType,
			Message:   fmt.Sprintf("Port %s %s", ev.Port.String(), ev.Type),
			Summary:   fmt.Sprintf("portwatch: port %s %s", ev.Port.String(), ev.Type),
			Payload: map[string]interface{}{
				"port":     ev.Port.Number,
				"protocol": ev.Port.Proto,
				"event":    string(ev.Type),
			},
		}

		body, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("zenduty: marshal: %w", err)
		}

		url := endpoint + z.integrationKey + "/"
		resp, err := z.client.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("zenduty: post: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("zenduty: unexpected status %d", resp.StatusCode)
		}
	}

	return nil
}
