package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func slackEvent(kind string) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Address: "0.0.0.0:8080", Proto: "tcp"},
	}
}

func TestSlackNotifierSuccess(t *testing.T) {
	var received slackPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL, ts.Client())
	if err := n.Notify(slackEvent("opened")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text == "" {
		t.Error("expected non-empty slack message text")
	}
}

func TestSlackNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL, ts.Client())
	err := n.Notify(slackEvent("closed"))
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestSlackNotifierBadURL(t *testing.T) {
	n := NewSlackNotifier("http://127.0.0.1:0/no-such-endpoint", nil)
	err := n.Notify(slackEvent("opened"))
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}
