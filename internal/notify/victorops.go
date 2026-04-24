package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// VictorOpsNotifier sends alerts to VictorOps (Splunk On-Call) via the REST endpoint.
type VictorOpsNotifier struct {
	endpointURL string
	routingKey  string
	client      *http.Client
}

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	Timestamp         int64  `json:"timestamp"`
}

// NewVictorOpsNotifier creates a notifier that posts to the given VictorOps REST endpoint URL.
// routingKey is appended to the URL path as required by the VictorOps API.
func NewVictorOpsNotifier(endpointURL, routingKey string) *VictorOpsNotifier {
	return &VictorOpsNotifier{
		endpointURL: endpointURL,
		routingKey:  routingKey,
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends each event as a separate VictorOps alert.
func (v *VictorOpsNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	url := fmt.Sprintf("%s/%s", v.endpointURL, v.routingKey)
	for _, e := range events {
		msgType := "INFO"
		if e.Kind == alert.Opened {
			msgType = "CRITICAL"
		}
		payload := victorOpsPayload{
			MessageType:       msgType,
			EntityID:          fmt.Sprintf("portwatch-%s", e.Port),
			EntityDisplayName: fmt.Sprintf("Port %s", e.Port),
			StateMessage:      fmt.Sprintf("Port %s %s", e.Port, e.Kind),
			Timestamp:         time.Now().Unix(),
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("victorops: marshal: %w", err)
		}
		resp, err := v.client.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("victorops: post: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("victorops: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}
