//go:build integration
// +build integration

package notify_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/scanner"
)

// TestProwlLiveNotify sends a real notification via the Prowl API.
// Requires PROWL_API_KEY environment variable to be set.
// Run with: go test -tags integration ./internal/notify/ -run TestProwlLiveNotify
func TestProwlLiveNotify(t *testing.T) {
	apiKey := os.Getenv("PROWL_API_KEY")
	if apiKey == "" {
		t.Skip("PROWL_API_KEY not set; skipping live integration test")
	}

	n := notify.NewProwlNotifier(apiKey, "portwatch-integration-test", 0)

	events := []alert.Event{
		{
			Kind: alert.Opened,
			Port: scanner.Port{Number: 8080, Protocol: "tcp"},
		},
		{
			Kind: alert.Closed,
			Port: scanner.Port{Number: 22, Protocol: "tcp"},
		},
	}

	if err := n.Notify(events); err != nil {
		t.Fatalf("live prowl notify failed: %v", err)
	}
	t.Log("Prowl notification sent successfully")
}
