//go:build integration
// +build integration

package notify

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// TestZendutyLiveNotify sends a real alert to Zenduty.
// Set ZENDUTY_INTEGRATION_KEY to run this test.
//
//	go test -tags integration -run TestZendutyLiveNotify ./internal/notify/
func TestZendutyLiveNotify(t *testing.T) {
	key := os.Getenv("ZENDUTY_INTEGRATION_KEY")
	if key == "" {
		t.Skip("ZENDUTY_INTEGRATION_KEY not set")
	}

	n := NewZendutyNotifier(key)

	events := []alert.Event{
		{
			Type: alert.Opened,
			Port: scanner.Port{Number: 9999, Proto: "tcp"},
		},
	}

	if err := n.Notify(events); err != nil {
		t.Fatalf("live notify failed: %v", err)
	}
}
