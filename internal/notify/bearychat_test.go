package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func bearyChatEvent(kind alert.EventKind) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: 8080, Protocol: "tcp"},
	}
}

func TestBearyChat_Success(t *testing.T) {
	var received bearyChatPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewBearyChat(srv.URL)
	events := []alert.Event{bearyChatEvent(alert.Opened), bearyChatEvent(alert.Closed)}

	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text == "" {
		t.Error("expected non-empty text in payload")
	}
	if len(received.Attachments) != 2 {
		t.Errorf("expected 2 attachments, got %d", len(received.Attachments))
	}
}

func TestBearyChat_Non2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := NewBearyChat(srv.URL)
	err := n.Notify([]alert.Event{bearyChatEvent(alert.Opened)})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestBearyChat_BadURL(t *testing.T) {
	n := NewBearyChat("http://127.0.0.1:0")
	err := n.Notify([]alert.Event{bearyChatEvent(alert.Opened)})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestBearyChat_NoEventsSkips(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewBearyChat(srv.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected HTTP call to be skipped when no events")
	}
}
