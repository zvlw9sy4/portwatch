//go:build integration
// +build integration

package notify

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// TestNewRelicLiveNotify sends a real event to New Relic.
//
// Run with:
//
//	NEW_RELIC_LICENSE_KEY=<key> go test -tags integration -run TestNewRelicLiveNotify ./internal/notify/
func TestNewRelicLiveNotify(t *testing.T) {
	key := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if key == "" {
		t.Skip("NEW_RELIC_LICENSE_KEY not set")
	}

	n := NewNewRelicNotifier(key, "")
	events := []alert.Event{
		{
			Port: scanner.Port{Number: 8080, Proto: "tcp"},
			Kind: alert.Opened,
		},
	}

	if err := n.Notify(events); err != nil {
		t.Fatalf("live New Relic notify failed: %v", err)
	}
	t.Log("event delivered to New Relic successfully")
}
