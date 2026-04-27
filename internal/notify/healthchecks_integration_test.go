//go:build integration
// +build integration

package notify

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// TestHealthChecksLiveNotify posts a real event to a healthchecks.io endpoint.
// Set HEALTHCHECKS_URL to a valid ping URL before running.
//
//	HEALTHCHECKS_URL=https://hc-ping.com/<uuid>/fail go test -tags integration -run TestHealthChecksLiveNotify ./internal/notify/
func TestHealthChecksLiveNotify(t *testing.T) {
	url := os.Getenv("HEALTHCHECKS_URL")
	if url == "" {
		t.Skip("HEALTHCHECKS_URL not set")
	}
	n := NewHealthChecksNotifier(url)
	events := []alert.Event{
		{
			Kind: alert.EventOpened,
			Port: scanner.Port{Number: 9999, Proto: "tcp"},
		},
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("live notify failed: %v", err)
	}
}
