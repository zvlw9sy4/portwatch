package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func gotifyEvent() alert.Event {
	return alert.Event{
		Kind: alert.Opened,
		Port: scanner.Port{Number: 9090, Protocol: "tcp"},
	}
}

func TestGotifyNotifierSuccess(t *testing.T) {
	var received gotifyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "testtoken", 5, ts.Client())
	if err := n.Notify(gotifyEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Priority != 5 {
		t.Errorf("priority = %d, want 5", received.Priority)
	}
	if received.Title == "" {
		t.Error("expected non-empty title")
	}
}

func TestGotifyNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "bad", 5, ts.Client())
	if err := n.Notify(gotifyEvent()); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestGotifyNotifierBadURL(t *testing.T) {
	n := NewGotifyNotifier("http://127.0.0.1:0", "tok", 5, nil)
	if err := n.Notify(gotifyEvent()); err == nil {
		t.Error("expected error for unreachable server")
	}
}
