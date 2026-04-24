package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func voEvent(kind alert.EventKind) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: 9090, Proto: "tcp"},
	}
}

func TestVictorOpsNotifierSuccess(t *testing.T) {
	var received victorOpsPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewVictorOpsNotifier(ts.URL, "mykey")
	if err := n.Notify([]alert.Event{voEvent(alert.Opened)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.MessageType != "CRITICAL" {
		t.Errorf("expected CRITICAL, got %s", received.MessageType)
	}
	if received.EntityID != "portwatch-9090/tcp" {
		t.Errorf("unexpected entity_id: %s", received.EntityID)
	}
}

func TestVictorOpsNotifierClosedIsInfo(t *testing.T) {
	var received victorOpsPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received) //nolint:errcheck
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewVictorOpsNotifier(ts.URL, "mykey")
	if err := n.Notify([]alert.Event{voEvent(alert.Closed)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.MessageType != "INFO" {
		t.Errorf("expected INFO, got %s", received.MessageType)
	}
}

func TestVictorOpsNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewVictorOpsNotifier(ts.URL, "mykey")
	err := n.Notify([]alert.Event{voEvent(alert.Opened)})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestVictorOpsNotifierBadURL(t *testing.T) {
	n := NewVictorOpsNotifier("http://127.0.0.1:1", "key")
	err := n.Notify([]alert.Event{voEvent(alert.Opened)})
	if err == nil {
		t.Fatal("expected error for unreachable host")
	}
}

func TestVictorOpsNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewVictorOpsNotifier(ts.URL, "mykey")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}
