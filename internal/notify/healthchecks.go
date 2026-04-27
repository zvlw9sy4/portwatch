package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// HealthChecksNotifier sends port change events to a healthchecks.io-compatible endpoint.
// Each alert posts a failure signal; when no events occur the daemon can ping the success URL.
type HealthChecksNotifier struct {
	pingURL string
	client  *http.Client
}

type healthChecksPayload struct {
	Event   string `json:"event"`
	Port    int    `json:"port"`
	Proto   string `json:"proto"`
	Message string `json:"message"`
}

// NewHealthChecksNotifier creates a notifier that posts to the given ping URL.
// Use the UUID-based URL from healthchecks.io, e.g. https://hc-ping.com/<uuid>/fail.
func NewHealthChecksNotifier(pingURL string) *HealthChecksNotifier {
	return &HealthChecksNotifier{
		pingURL: pingURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends one HTTP POST per event to the configured ping URL.
func (h *HealthChecksNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	var errs []error
	for _, ev := range events {
		payload := healthChecksPayload{
			Event:   string(ev.Kind),
			Port:    ev.Port.Number,
			Proto:   ev.Port.Proto,
			Message: fmt.Sprintf("port %d/%s %s", ev.Port.Number, ev.Port.Proto, ev.Kind),
		}
		body, err := json.Marshal(payload)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		resp, err := h.client.Post(h.pingURL, "application/json", bytes.NewReader(body))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			errs = append(errs, fmt.Errorf("healthchecks: unexpected status %d", resp.StatusCode))
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
