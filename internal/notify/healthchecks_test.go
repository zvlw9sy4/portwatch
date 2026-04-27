package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func hcEvent(kind alert.EventKind, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: port, Proto: "tcp"},
	}
}

func TestHealthChecksNotifierSuccess(t *testing.T) {
	var received []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&p)
		received = append(received, p)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewHealthChecksNotifier(ts.URL)
	events := []alert.Event{hcEvent(alert.EventOpened, 8080), hcEvent(alert.EventClosed, 9090)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(received))
	}
	if received[0]["port"].(float64) != 8080 {
		t.Errorf("expected port 8080, got %v", received[0]["port"])
	}
}

func TestHealthChecksNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewHealthChecksNotifier(ts.URL)
	err := n.Notify([]alert.Event{hcEvent(alert.EventOpened, 443)})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestHealthChecksNotifierBadURL(t *testing.T) {
	n := NewHealthChecksNotifier("http://127.0.0.1:0")
	err := n.Notify([]alert.Event{hcEvent(alert.EventOpened, 22)})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestHealthChecksNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewHealthChecksNotifier(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}
