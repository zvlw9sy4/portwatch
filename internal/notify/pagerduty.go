package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyNotifier sends alerts to PagerDuty via the Events API v2.
type PagerDutyNotifier struct {
	integrationKey string
	client         *http.Client
	endpoint       string
}

// NewPagerDutyNotifier creates a notifier that posts to PagerDuty.
func NewPagerDutyNotifier(integrationKey string) *PagerDutyNotifier {
	return &PagerDutyNotifier{
		integrationKey: integrationKey,
		client:         &http.Client{},
		endpoint:       pagerDutyEventURL,
	}
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary  string `json:"summary"`
	Severity string `json:"severity"`
	Source   string `json:"source"`
}

// Notify sends a single alert event to PagerDuty.
func (p *PagerDutyNotifier) Notify(e alert.Event) error {
	body := pdPayload{
		RoutingKey:  p.integrationKey,
		EventAction: "trigger",
		Payload: pdDetails{
			Summary:  fmt.Sprintf("portwatch: port %d/%s %s", e.Port.Number, e.Port.Protocol, e.Kind),
			Severity: "warning",
			Source:   "portwatch",
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal: %w", err)
	}
	resp, err := p.client.Post(p.endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
