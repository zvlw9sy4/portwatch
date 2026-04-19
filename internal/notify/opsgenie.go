package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// OpsGenieNotifier sends alerts to OpsGenie via its REST API.
type OpsGenieNotifier struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

type opsGeniePayload struct {
	Message     string            `json:"message"`
	Description string            `json:"description"`
	Priority    string            `json:"priority"`
	Tags        []string          `json:"tags"`
	Details     map[string]string `json:"details"`
}

// NewOpsGenieNotifier creates an OpsGenieNotifier with the given API key.
func NewOpsGenieNotifier(apiKey string) *OpsGenieNotifier {
	return &OpsGenieNotifier{
		apiKey: apiKey,
		apiURL: "https://api.opsgenie.com/v2/alerts",
		client: &http.Client{},
	}
}

func (o *OpsGenieNotifier) Notify(e alert.Event) error {
	payload := opsGeniePayload{
		Message:     fmt.Sprintf("portwatch: port %s %s", e.Port, e.Kind),
		Description: fmt.Sprintf("Port %d/%s was %s", e.Port.Number, e.Port.Protocol, e.Kind),
		Priority:    "P2",
		Tags:        []string{"portwatch", string(e.Kind)},
		Details: map[string]string{
			"port":     fmt.Sprintf("%d", e.Port.Number),
			"protocol": e.Port.Protocol,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, o.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}
