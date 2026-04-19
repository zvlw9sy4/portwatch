package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func discordEvent(kind alert.EventKind) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: 9200, Proto: "tcp"},
	}
}

func TestDiscordNotifierSuccess(t *testing.T) {
	var got discordPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL)
	if err := n.Notify(discordEvent(alert.Opened)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Content == "" {
		t.Error("expected non-empty content")
	}
}

func TestDiscordNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL)
	if err := n.Notify(discordEvent(alert.Closed)); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestDiscordNotifierBadURL(t *testing.T) {
	n := NewDiscordNotifier("http://127.0.0.1:0")
	if err := n.Notify(discordEvent(alert.Opened)); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
