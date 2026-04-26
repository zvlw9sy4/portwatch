package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// CustomEventNotifier sends port change events to a generic JSON webhook
// endpoint with a configurable payload template via field mappings.
type CustomEventNotifier struct {
	url     string
	apiKey  string
	client  *http.Client
}

type customEventPayload struct {
	EventType string `json:"event_type"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
	Message   string `json:"message"`
}

// NewCustomEventNotifier creates a CustomEventNotifier that POSTs events to url,
// attaching apiKey as a Bearer token when non-empty.
func NewCustomEventNotifier(url, apiKey string) *CustomEventNotifier {
	return &CustomEventNotifier{
		url:    url,
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// Notify sends each alert.Event as a separate JSON POST request.
func (n *CustomEventNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	var errs []error
	for _, e := range events {
		if err := n.send(e); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("customevent: %d error(s): %v", len(errs), errs[0])
	}
	return nil
}

func (n *CustomEventNotifier) send(e alert.Event) error {
	payload := customEventPayload{
		EventType: string(e.Kind),
		Port:      e.Port.Number,
		Protocol:  e.Port.Protocol,
		Message:   fmt.Sprintf("port %s %s", e.Port.String(), e.Kind),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("customevent: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, n.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("customevent: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if n.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+n.apiKey)
	}
	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("customevent: request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("customevent: unexpected status %d", resp.StatusCode)
	}
	return nil
}
