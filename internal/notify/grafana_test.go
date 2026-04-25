package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func grafanaEvent() alert.Event {
	return alert.Event{
		Kind: "opened",
		Port: scanner.Port{Number: 9090, Proto: "tcp"},
	}
}

func TestGrafanaNotifierSuccess(t *testing.T) {
	var received grafanaAnnotation
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer testkey" {
			t.Errorf("missing or wrong Authorization header: %s", r.Header.Get("Authorization"))
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGrafanaNotifier(ts.URL, "testkey", []string{"portwatch", "prod"})
	if err := n.Notify([]alert.Event{grafanaEvent()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text == "" {
		t.Error("expected annotation text to be set")
	}
	if len(received.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(received.Tags))
	}
}

func TestGrafanaNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewGrafanaNotifier(ts.URL, "bad", nil)
	if err := n.Notify([]alert.Event{grafanaEvent()}); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestGrafanaNotifierBadURL(t *testing.T) {
	n := NewGrafanaNotifier("://bad-url", "key", nil)
	if err := n.Notify([]alert.Event{grafanaEvent()}); err == nil {
		t.Error("expected error on bad URL")
	}
}

func TestGrafanaNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewGrafanaNotifier(ts.URL, "key", nil)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}
