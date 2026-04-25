//go:build integration
// +build integration

package notify_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/scanner"
)

// TestOpsGenieLiveNotify sends a real alert to OpsGenie.
// Requires OPSGENIE_API_KEY to be set in the environment.
// Run with: go test -tags integration -run TestOpsGenieLiveNotify ./internal/notify/
func TestOpsGenieLiveNotify(t *testing.T) {
	apiKey := os.Getenv("OPSGENIE_API_KEY")
	if apiKey == "" {
		t.Skip("OPSGENIE_API_KEY not set")
	}

	n := notify.NewOpsGenieNotifier(notify.OpsGenieConfig{
		APIKey:   apiKey,
		Region:   "us",
		Priority: "P5",
		Tags:     []string{"portwatch", "integration-test"},
	})

	evt := ogEvent(scanner.Port{Number: 9999, Protocol: "tcp"})
	if err := n.Notify([]interface{}{evt}); err != nil {
		t.Fatalf("live notify failed: %v", err)
	}
}
