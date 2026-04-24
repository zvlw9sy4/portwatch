package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const pushoverAPIURL = "https://api.pushover.net/1/messages.json"

// PushoverNotifier sends alerts via the Pushover notification service.
type PushoverNotifier struct {
	token   string
	user    string
	apiURL  string
	client  *http.Client
}

type pushoverPayload struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// NewPushoverNotifier creates a notifier that delivers messages via Pushover.
// token is the application API token and user is the user/group key.
func NewPushoverNotifier(token, user string) *PushoverNotifier {
	return &PushoverNotifier{
		token:  token,
		user:   user,
		apiURL: pushoverAPIURL,
		client: &http.Client{},
	}
}

// Notify sends each event as a separate Pushover message.
func (p *PushoverNotifier) Notify(events []alert.Event) error {
	for _, ev := range events {
		if err := p.send(ev); err != nil {
			return err
		}
	}
	return nil
}

func (p *PushoverNotifier) send(ev alert.Event) error {
	body, err := json.Marshal(pushoverPayload{
		Token:   p.token,
		User:    p.user,
		Title:   "portwatch: port " + ev.Kind,
		Message: fmt.Sprintf("Port %s/%s %s", ev.Port.Proto, ev.Port.String(), ev.Kind),
	})
	if err != nil {
		return fmt.Errorf("pushover: marshal: %w", err)
	}

	resp, err := p.client.Post(p.apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pushover: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pushover: unexpected status %d", resp.StatusCode)
	}
	return nil
}
