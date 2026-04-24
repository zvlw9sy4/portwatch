package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// SignaldNotifier sends alerts via the signald HTTP REST API (signal-cli-rest-api).
type SignaldNotifier struct {
	baseURL    string
	sender     string
	recipients []string
	client     *http.Client
}

type signaldPayload struct {
	Message    string   `json:"message"`
	Recipients []string `json:"recipients"`
	Number     string   `json:"number"`
}

// NewSignaldNotifier creates a SignaldNotifier that posts to the given
// signal-cli-rest-api base URL (e.g. "http://localhost:8080").
func NewSignaldNotifier(baseURL, sender string, recipients []string) *SignaldNotifier {
	return &SignaldNotifier{
		baseURL:    baseURL,
		sender:     sender,
		recipients: recipients,
		client:     &http.Client{},
	}
}

// Notify sends all events as a single Signal message.
func (s *SignaldNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, e := range events {
		fmt.Fprintf(&buf, "[%s] port %s\n", e.Kind, e.Port)
	}

	payload := signaldPayload{
		Message:    buf.String(),
		Recipients: s.recipients,
		Number:     s.sender,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signald: marshal payload: %w", err)
	}

	url := s.baseURL + "/v2/send"
	resp, err := s.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signald: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signald: unexpected status %d", resp.StatusCode)
	}
	return nil
}
