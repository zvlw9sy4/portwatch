package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func mmEvent(kind alert.EventKind, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: port, Proto: "tcp"},
	}
}

func TestMattermostNotifierSuccess(t *testing.T) {
	var received mattermostPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewMattermostNotifier(ts.URL, "#alerts")
	events := []alert.Event{
		mmEvent(alert.Opened, 8080),
		mmEvent(alert.Closed, 22),
	}

	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Channel != "#alerts" {
		t.Errorf("channel = %q, want #alerts", received.Channel)
	}
	if received.Text == "" {
		t.Error("expected non-empty text")
	}
}

func TestMattermostNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewMattermostNotifier(ts.URL, "")
	err := n.Notify([]alert.Event{mmEvent(alert.Opened, 9090)})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestMattermostNotifierBadURL(t *testing.T) {
	n := NewMattermostNotifier("http://127.0.0.1:0/no-server", "")
	err := n.Notify([]alert.Event{mmEvent(alert.Opened, 443)})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestMattermostNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewMattermostNotifier(ts.URL, "")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected HTTP call to be skipped for empty event list")
	}
}
