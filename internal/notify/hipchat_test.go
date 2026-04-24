package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func hipchatEvent() alert.Event {
	return alert.Event{
		Kind: alert.Opened,
		Port: scanner.Port{Number: 8080, Protocol: "tcp"},
	}
}

func TestHipChatNotifierSuccess(t *testing.T) {
	var received hipChatPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("unexpected auth header: %s", auth)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	n := NewHipChatNotifier("42", "test-token", srv.URL)
	if err := n.Notify([]alert.Event{hipchatEvent()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Color != "green" {
		t.Errorf("expected green color for opened event, got %s", received.Color)
	}
	if !received.Notify {
		t.Error("expected notify=true")
	}
}

func TestHipChatNotifierNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	n := NewHipChatNotifier("42", "bad-token", srv.URL)
	if err := n.Notify([]alert.Event{hipchatEvent()}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestHipChatNotifierBadURL(t *testing.T) {
	n := NewHipChatNotifier("42", "tok", "://bad-url")
	if err := n.Notify([]alert.Event{hipchatEvent()}); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestHipChatNotifierNoEventsSkips(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer srv.Close()

	n := NewHipChatNotifier("42", "tok", srv.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}

func TestHipChatNotifierClosedEventColor(t *testing.T) {
	var received hipChatPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	closedEv := alert.Event{
		Kind: alert.Closed,
		Port: scanner.Port{Number: 443, Protocol: "tcp"},
	}

	n := NewHipChatNotifier("42", "tok", srv.URL)
	if err := n.Notify([]alert.Event{closedEv}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Color != "yellow" {
		t.Errorf("expected yellow for closed event, got %s", received.Color)
	}
}
