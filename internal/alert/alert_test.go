package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(n int) scanner.Port {
	return scanner.Port{Number: n, Protocol: "tcp"}
}

func TestBuildEventsOpened(t *testing.T) {
	events := alert.BuildEvents([]scanner.Port{makePort(8080)}, nil)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != alert.LevelAlert {
		t.Errorf("expected ALERT level, got %s", events[0].Level)
	}
}

func TestBuildEventsClosed(t *testing.T) {
	events := alert.BuildEvents(nil, []scanner.Port{makePort(22)})
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != alert.LevelWarn {
		t.Errorf("expected WARN level, got %s", events[0].Level)
	}
}

func TestBuildEventsBoth(t *testing.T) {
	events := alert.BuildEvents(
		[]scanner.Port{makePort(443)},
		[]scanner.Port{makePort(80), makePort(8443)},
	)
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
}

func TestLogNotifierOutput(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)
	events := alert.BuildEvents([]scanner.Port{makePort(9090)}, nil)
	for _, e := range events {
		if err := n.Notify(e); err != nil {
			t.Fatalf("Notify error: %v", err)
		}
	}
	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "9090") {
		t.Errorf("expected port 9090 in output, got: %s", out)
	}
}

func TestLogNotifierDefaultWriter(t *testing.T) {
	n := alert.NewLogNotifier(nil)
	if n.Out == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}
