package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

const prowlAPIURL = "https://api.prowlapp.com/publicapi/add"

// ProwlNotifier sends push notifications to iOS devices via the Prowl API.
type ProwlNotifier struct {
	apiKey      string
	appName     string
	priority    int
	httpClient  *http.Client
}

// NewProwlNotifier creates a Prowl notifier.
// priority ranges from -2 (very low) to 2 (emergency).
func NewProwlNotifier(apiKey, appName string, priority int) *ProwlNotifier {
	if appName == "" {
		appName = "portwatch"
	}
	if priority < -2 || priority > 2 {
		priority = 0
	}
	return &ProwlNotifier{
		apiKey:     apiKey,
		appName:    appName,
		priority:   priority,
		httpClient: &http.Client{},
	}
}

// Notify sends port change events to Prowl.
func (p *ProwlNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var lines []string
	for _, e := range events {
		lines = append(lines, fmt.Sprintf("%s: %s", e.Kind, e.Port))
	}
	description := strings.Join(lines, "\n")

	params := url.Values{}
	params.Set("apikey", p.apiKey)
	params.Set("application", p.appName)
	params.Set("event", "Port Change Detected")
	params.Set("description", description)
	params.Set("priority", fmt.Sprintf("%d", p.priority))

	resp, err := p.httpClient.PostForm(prowlAPIURL, params)
	if err != nil {
		return fmt.Errorf("prowl: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("prowl: unexpected status %d", resp.StatusCode)
	}
	return nil
}
