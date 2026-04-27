package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

const pushbulletAPIURL = "https://api.pushbullet.com/v2/pushes"

// PushbulletNotifier sends port-change alerts via the Pushbullet API.
// It requires a valid API access token from https://www.pushbullet.com/#settings/account.
type PushbulletNotifier struct {
	token  string
	client *http.Client
}

// NewPushbulletNotifier creates a PushbulletNotifier that authenticates
// with the given API access token.
func NewPushbulletNotifier(token string) *PushbulletNotifier {
	return &PushbulletNotifier{
		token:  token,
		client: &http.Client{},
	}
}

// Notify sends a Pushbullet note push for each alert event.
// Events with no changes are silently skipped.
func (p *PushbulletNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	title, body := buildPushbulletMessage(events)

	payload := map[string]string{
		"type":  "note",
		"title": title,
		"body":  body,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pushbullet: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, pushbulletAPIURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pushbullet: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("pushbullet: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pushbullet: unexpected status %d", resp.StatusCode)
	}

	return nil
}

// buildPushbulletMessage returns a title and body summarising the given events.
func buildPushbulletMessage(events []alert.Event) (title, body string) {
	var opened, closed []string
	for _, e := range events {
		desc := fmt.Sprintf("%s/%d", e.Port.Protocol, e.Port.Number)
		switch e.Kind {
		case alert.Opened:
			opened = append(opened, desc)
		case alert.Closed:
			closed = append(closed, desc)
		}
	}

	var parts []string
	if len(opened) > 0 {
		parts = append(parts, fmt.Sprintf("Opened: %s", strings.Join(opened, ", ")))
	}
	if len(closed) > 0 {
		parts = append(parts, fmt.Sprintf("Closed: %s", strings.Join(closed, ", ")))
	}

	title = fmt.Sprintf("portwatch: %d port change(s) detected", len(events))
	body = strings.Join(parts, "\n")
	return title, body
}
