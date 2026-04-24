package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// BearyChat incoming webhook payload.
type bearyChatPayload struct {
	Text        string             `json:"text"`
	Attachments []bearyChatAttach  `json:"attachments,omitempty"`
}

type bearyChatAttach struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}

// BearyChat notifier sends port-change events to a BearyChat incoming webhook.
type BearyChat struct {
	webhookURL string
	client     *http.Client
}

// NewBearyChat returns a Notifier that posts to the given BearyChat webhook URL.
func NewBearyChat(webhookURL string) *BearyChat {
	return &BearyChat{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Notify sends all events to BearyChat as a single message.
func (b *BearyChat) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	attachments := make([]bearyChatAttach, 0, len(events))
	for _, e := range events {
		color := "#36a64f"
		if e.Kind == alert.Closed {
			color = "#e01e5a"
		}
		attachments = append(attachments, bearyChatAttach{
			Text:  fmt.Sprintf("%s %s", e.Kind, e.Port),
			Color: color,
		})
	}

	payload := bearyChatPayload{
		Text:        "portwatch: port changes detected",
		Attachments: attachments,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("bearychat: marshal: %w", err)
	}

	resp, err := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("bearychat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bearychat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
