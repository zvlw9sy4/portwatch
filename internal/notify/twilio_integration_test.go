//go:build integration
// +build integration

package notify

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// TestTwilioLiveNotify sends a real SMS via Twilio.
// Requires environment variables:
//
//	TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_FROM, TWILIO_TO
//
// Run with: go test -tags integration -run TestTwilioLiveNotify ./internal/notify/
func TestTwilioLiveNotify(t *testing.T) {
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	token := os.Getenv("TWILIO_AUTH_TOKEN")
	from := os.Getenv("TWILIO_FROM")
	to := os.Getenv("TWILIO_TO")

	if sid == "" || token == "" || from == "" || to == "" {
		t.Skip("twilio integration env vars not set")
	}

	n := NewTwilioNotifier(sid, token, from, to)

	events := []alert.Event{
		{Kind: alert.Opened, Port: scanner.Port{Number: 8080, Protocol: "tcp"}},
		{Kind: alert.Closed, Port: scanner.Port{Number: 9090, Protocol: "tcp"}},
	}

	if err := n.Notify(events); err != nil {
		t.Fatalf("live twilio notify failed: %v", err)
	}
}
