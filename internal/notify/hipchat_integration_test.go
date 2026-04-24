//go:build integration
// +build integration

package notify

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// TestHipChatLiveNotify sends a real notification to HipChat.
// Run with:
//
//	HIPCHAT_ROOM=<room> HIPCHAT_TOKEN=<token> HIPCHAT_URL=<url> \
//	  go test -tags integration ./internal/notify/ -run TestHipChatLiveNotify
func TestHipChatLiveNotify(t *testing.T) {
	room := os.Getenv("HIPCHAT_ROOM")
	token := os.Getenv("HIPCHAT_TOKEN")
	baseURL := os.Getenv("HIPCHAT_URL")

	if room == "" || token == "" || baseURL == "" {
		t.Skip("HIPCHAT_ROOM, HIPCHAT_TOKEN, HIPCHAT_URL not set")
	}

	n := NewHipChatNotifier(room, token, baseURL)

	events := []alert.Event{
		{
			Kind: alert.Opened,
			Port: scanner.Port{Number: 9999, Protocol: "tcp"},
		},
	}

	if err := n.Notify(events); err != nil {
		t.Fatalf("live HipChat notify failed: %v", err)
	}
}
