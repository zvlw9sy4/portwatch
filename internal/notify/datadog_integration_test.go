//go:build integration
// +build integration

package notify

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// TestDatadogLiveNotify sends a real event to Datadog.
// Run with: DD_API_KEY=<key> go test -tags integration ./internal/notify/
func TestDatadogLiveNotify(t *testing.T) {
	apiKey := os.Getenv("DD_API_KEY")
	if apiKey == "" {
		t.Skip("DD_API_KEY not set")
	}

	n := NewDatadogNotifier(apiKey, "")
	events := []alert.Event{
		{
			Kind: alert.Opened,
			Port: scanner.Port{Number: 9999, Proto: "tcp"},
		},
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("live notify failed: %v", err)
	}
}
