//go:build integration
// +build integration

package notify

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// TestSplunkLiveNotify sends a real event to a Splunk HEC endpoint.
// Set SPLUNK_HEC_URL and SPLUNK_HEC_TOKEN env vars before running.
//
//	SPLUNK_HEC_URL=https://splunk:8088/services/collector \
//	SPLUNK_HEC_TOKEN=abc123 \
//	go test -tags integration -run TestSplunkLiveNotify ./internal/notify/
func TestSplunkLiveNotify(t *testing.T) {
	url := os.Getenv("SPLUNK_HEC_URL")
	token := os.Getenv("SPLUNK_HEC_TOKEN")
	if url == "" || token == "" {
		t.Skip("SPLUNK_HEC_URL or SPLUNK_HEC_TOKEN not set")
	}

	n := NewSplunkNotifier(url, token, "portwatch-integration-test")
	events := []alert.Event{
		{
			Kind: "opened",
			Port: scanner.Port{Number: 8080, Proto: "tcp"},
			Time: time.Now(),
		},
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("live Splunk notify failed: %v", err)
	}
}
