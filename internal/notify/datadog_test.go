package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func datadogEvent() alert.Event {
	return alert.Event{
		Kind: alert.Opened,
		Port: scanner.Port{Number: 9200, Proto: "tcp"},
	}
}

func TestDatadogNotifierSuccess(t *testing.T) {
	var received []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		received = append(received, payload)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewDatadogNotifier("test-key", ts.URL)
	if err := n.Notify([]alert.Event{datadogEvent()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 1 {
		t.Fatalf("expected 1 request, got %d", len(received))
	}
	if received[0]["alert_type"] != "warning" {
		t.Errorf("expected alert_type=warning, got %v", received[0]["alert_type"])
	}
}

func TestDatadogNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewDatadogNotifier("bad-key", ts.URL)
	if err := n.Notify([]alert.Event{datadogEvent()}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestDatadogNotifierBadURL(t *testing.T) {
	n := NewDatadogNotifier("key", "http://127.0.0.1:0")
	if err := n.Notify([]alert.Event{datadogEvent()}); err == nil {
		t.Fatal("expected error for unreachable host")
	}
}

func TestDatadogNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewDatadogNotifier("key", ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}
