package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func rcEvent() alert.Event {
	return alert.Event{
		Kind: "opened",
		Port: scanner.Port{Number: 9200, Proto: "tcp"},
	}
}

func TestRocketChatNotifierSuccess(t *testing.T) {
	var received rocketChatPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewRocketChatNotifier(ts.URL)
	if err := n.Notify([]alert.Event{rcEvent()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text == "" {
		t.Error("expected non-empty text payload")
	}
	if want := "opened"; !containsStr(received.Text, want) {
		t.Errorf("payload text %q missing %q", received.Text, want)
	}
}

func TestRocketChatNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewRocketChatNotifier(ts.URL)
	if err := n.Notify([]alert.Event{rcEvent()}); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestRocketChatNotifierBadURL(t *testing.T) {
	n := NewRocketChatNotifier("http://127.0.0.1:0/no-server")
	if err := n.Notify([]alert.Event{rcEvent()}); err == nil {
		t.Error("expected error for unreachable URL")
	}
}

func TestRocketChatNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewRocketChatNotifier(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}

// containsStr is a small helper shared across notify tests.
func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
