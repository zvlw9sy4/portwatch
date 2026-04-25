package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// GrafanaNotifier sends port-change events to a Grafana annotation endpoint.
type GrafanaNotifier struct {
	url    string
	apiKey string
	tags   []string
	client *http.Client
}

type grafanaAnnotation struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

// NewGrafanaNotifier creates a notifier that posts annotations to Grafana.
// url should be the full annotations API endpoint, e.g.
// "http://grafana:3000/api/annotations".
func NewGrafanaNotifier(url, apiKey string, tags []string) *GrafanaNotifier {
	if tags == nil {
		tags = []string{"portwatch"}
	}
	return &GrafanaNotifier{
		url:    url,
		apiKey: apiKey,
		tags:   tags,
		client: &http.Client{},
	}
}

// Notify posts one Grafana annotation per event.
func (g *GrafanaNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	var errs []error
	for _, ev := range events {
		text := fmt.Sprintf("[portwatch] port %s %s", ev.Port, ev.Kind)
		body, err := json.Marshal(grafanaAnnotation{Text: text, Tags: g.tags})
		if err != nil {
			errs = append(errs, err)
			continue
		}
		req, err := http.NewRequest(http.MethodPost, g.url, bytes.NewReader(body))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+g.apiKey)
		resp, err := g.client.Do(req)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			errs = append(errs, fmt.Errorf("grafana: unexpected status %d", resp.StatusCode))
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
