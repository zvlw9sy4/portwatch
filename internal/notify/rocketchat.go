package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// RocketChatNotifier sends port change events to a Rocket.Chat
// incoming webhook URL.
type RocketChatNotifier struct {
	webhookURL string
	client     *http.Client
}

type rocketChatPayload struct {
	Text string `json:"text"`
}

// NewRocketChatNotifier returns a RocketChatNotifier that posts
// messages to the given Rocket.Chat incoming webhook URL.
func NewRocketChatNotifier(webhookURL string) *RocketChatNotifier {
	return &RocketChatNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Notify sends all events to the configured Rocket.Chat webhook.
func (r *RocketChatNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("[portwatch] %s port %s/%d\n",
			e.Kind, e.Port.Proto, e.Port.Number))
	}

	payload := rocketChatPayload{Text: buf.String()}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rocketchat: marshal payload: %w", err)
	}

	resp, err := r.client.Post(r.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("rocketchat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rocketchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
