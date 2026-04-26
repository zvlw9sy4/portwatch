package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func googleChatEvent(kind alert.EventKind, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: port, Protocol: "tcp"},
	}
}

func TestGoogleChatNotifierSuccess(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	events := []alert.Event{
		googleChatEvent(alert.Opened, 8080),
		googleChatEvent(alert.Closed, 22),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] == "" {
		t.Error("expected non-empty text field in payload")
	}
}

func TestGoogleChatNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	err := n.Notify([]alert.Event{googleChatEvent(alert.Opened, 9090)})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestGoogleChatNotifierBadURL(t *testing.T) {
	n := NewGoogleChatNotifier("http://127.0.0.1:0/no-server")
	err := n.Notify([]alert.Event{googleChatEvent(alert.Opened, 443)})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestGoogleChatNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewGoogleChatNotifier(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}
