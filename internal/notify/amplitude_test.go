package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func amplitudeEvent() alert.Event {
	return alert.Event{
		Kind: "opened",
		Port: scanner.Port{Number: 9090, Proto: "tcp"},
	}
}

func TestAmplitudeNotifierSuccess(t *testing.T) {
	var received amplitudePayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewAmplitudeNotifier("test-key", srv.URL)
	if err := n.Notify([]alert.Event{amplitudeEvent()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.APIKey != "test-key" {
		t.Errorf("expected api_key %q, got %q", "test-key", received.APIKey)
	}
	if len(received.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(received.Events))
	}
	if received.Events[0].EventType != "port_opened" {
		t.Errorf("unexpected event_type %q", received.Events[0].EventType)
	}
	if received.Events[0].EventProps["port"] != float64(9090) {
		t.Errorf("unexpected port %v", received.Events[0].EventProps["port"])
	}
}

func TestAmplitudeNotifierNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := NewAmplitudeNotifier("key", srv.URL)
	if err := n.Notify([]alert.Event{amplitudeEvent()}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestAmplitudeNotifierBadURL(t *testing.T) {
	n := NewAmplitudeNotifier("key", "http://127.0.0.1:0")
	if err := n.Notify([]alert.Event{amplitudeEvent()}); err == nil {
		t.Fatal("expected error for unreachable host")
	}
}

func TestAmplitudeNotifierNoEventsSkips(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewAmplitudeNotifier("key", srv.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}
