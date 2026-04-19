package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func telegramEvent() alert.Event {
	return alert.Event{
		Kind: alert.Opened,
		Port: scanner.Port{Number: 8443, Protocol: "tcp"},
	}
}

func TestTelegramNotifierSuccess(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTelegramNotifier("tok123", "chat456")
	n.apiBase = ts.URL

	if err := n.Notify(telegramEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["chat_id"] != "chat456" {
		t.Errorf("chat_id = %q, want chat456", got["chat_id"])
	}
	if got["text"] == "" {
		t.Error("expected non-empty text")
	}
}

func TestTelegramNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewTelegramNotifier("tok", "chat")
	n.apiBase = ts.URL

	if err := n.Notify(telegramEvent()); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestTelegramNotifierBadURL(t *testing.T) {
	n := NewTelegramNotifier("tok", "chat")
	n.apiBase = "http://127.0.0.1:0"

	if err := n.Notify(telegramEvent()); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
