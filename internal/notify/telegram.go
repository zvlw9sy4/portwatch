package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const telegramAPIBase = "https://api.telegram.org"

// TelegramNotifier sends port-change alerts to a Telegram chat via Bot API.
type TelegramNotifier struct {
	token   string
	chatID  string
	apiBase string
	client  *http.Client
}

// NewTelegramNotifier returns a notifier that posts messages to the given chat.
func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		token:   token,
		chatID:  chatID,
		apiBase: telegramAPIBase,
		client:  &http.Client{},
	}
}

func (t *TelegramNotifier) Notify(e alert.Event) error {
	text := fmt.Sprintf("[portwatch] %s port %s/%d",
		e.Kind, e.Port.Protocol, e.Port.Number)

	payload, err := json.Marshal(map[string]string{
		"chat_id": t.chatID,
		"text":    text,
	})
	if err != nil {
		return fmt.Errorf("telegram: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", t.apiBase, t.token)
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("telegram: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}
	return nil
}
