package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// DingTalkNotifier sends port change alerts to a DingTalk webhook.
type DingTalkNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewDingTalkNotifier creates a new DingTalkNotifier that posts to the given webhook URL.
func NewDingTalkNotifier(webhookURL string) *DingTalkNotifier {
	return &DingTalkNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

type dingtalkPayload struct {
	MsgType string          `json:"msgtype"`
	Text    dingtalkText    `json:"text"`
	At      dingtalkAt      `json:"at"`
}

type dingtalkText struct {
	Content string `json:"content"`
}

type dingtalkAt struct {
	IsAtAll bool `json:"isAtAll"`
}

// Notify sends all events to DingTalk as a single text message.
func (d *DingTalkNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var buf bytes.Buffer
	buf.WriteString("portwatch alert:\n")
	for _, e := range events {
		buf.WriteString(fmt.Sprintf("  [%s] %s/%d\n", e.Kind, e.Port.Proto, e.Port.Number))
	}

	payload := dingtalkPayload{
		MsgType: "text",
		Text:    dingtalkText{Content: buf.String()},
		At:      dingtalkAt{IsAtAll: false},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("dingtalk: marshal payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("dingtalk: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("dingtalk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
