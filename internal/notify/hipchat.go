package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// HipChatNotifier sends port-change alerts to a HipChat room via the v2 API.
type HipChatNotifier struct {
	roomID  string
	token   string
	baseURL string
	client  *http.Client
}

type hipChatPayload struct {
	Message       string `json:"message"`
	MessageFormat string `json:"message_format"`
	Color         string `json:"color"`
	Notify        bool   `json:"notify"`
}

// NewHipChatNotifier creates a notifier that posts to the given HipChat room.
// baseURL should be the HipChat server root, e.g. "https://api.hipchat.com".
func NewHipChatNotifier(roomID, token, baseURL string) *HipChatNotifier {
	return &HipChatNotifier{
		roomID:  roomID,
		token:   token,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Notify sends each event as a separate HipChat room notification.
func (h *HipChatNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	for _, ev := range events {
		color := "green"
		if ev.Kind == alert.Closed {
			color = "yellow"
		}

		payload := hipChatPayload{
			Message:       fmt.Sprintf("[portwatch] %s", ev.String()),
			MessageFormat: "text",
			Color:         color,
			Notify:        true,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("hipchat: marshal: %w", err)
		}

		url := fmt.Sprintf("%s/v2/room/%s/notification", h.baseURL, h.roomID)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("hipchat: request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+h.token)

		resp, err := h.client.Do(req)
		if err != nil {
			return fmt.Errorf("hipchat: send: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("hipchat: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}
