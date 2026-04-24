//go:build integration
// +build integration

package notify

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// TestVictorOpsLiveNotify sends a real alert to VictorOps.
// Requires VICTOROPS_ENDPOINT and VICTOROPS_ROUTING_KEY environment variables.
// Run with: go test -tags integration ./internal/notify/
func TestVictorOpsLiveNotify(t *testing.T) {
	endpoint := os.Getenv("VICTOROPS_ENDPOINT")
	routingKey := os.Getenv("VICTOROPS_ROUTING_KEY")
	if endpoint == "" || routingKey == "" {
		t.Skip("VICTOROPS_ENDPOINT or VICTOROPS_ROUTING_KEY not set")
	}

	n := NewVictorOpsNotifier(endpoint, routingKey)
	events := []alert.Event{
		{
			Kind: alert.Opened,
			Port: scanner.Port{Number: 8080, Proto: "tcp"},
		},
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("live notify failed: %v", err)
	}
}
