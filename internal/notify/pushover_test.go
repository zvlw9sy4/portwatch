package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func pushoverEvent() alert.Event {
	return alert.Event{
		Kind: "opened",
		Port: scanner.Port{Number: 9090, Proto: "tcp"},
	}
}

func TestPushoverNotifierSuccess(t *testing.T) {
	var received pushoverPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewPushoverNotifier("app-token", "user-key")
	n.apiURL = srv.URL

	if err := n.Notify([]alert.Event{pushoverEvent()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Token != "app-token" {
		t.Errorf("expected token app-token, got %q", received.Token)
	}
	if received.User != "user-key" {
		t.Errorf("expected user user-key, got %q", received.User)
	}
	if received.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestPushoverNotifierNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	n := NewPushoverNotifier("tok", "usr")
	n.apiURL = srv.URL

	if err := n.Notify([]alert.Event{pushoverEvent()}); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestPushoverNotifierBadURL(t *testing.T) {
	n := NewPushoverNotifier("tok", "usr")
	n.apiURL = "http://127.0.0.1:0"

	if err := n.Notify([]alert.Event{pushoverEvent()}); err == nil {
		t.Error("expected error on bad URL")
	}
}
