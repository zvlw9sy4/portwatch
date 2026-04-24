package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func splunkEvent() alert.Event {
	return alert.Event{
		Kind: "opened",
		Port: scanner.Port{Number: 9200, Proto: "tcp"},
		Time: time.Unix(1700000000, 0),
	}
}

func TestSplunkNotifierSuccess(t *testing.T) {
	var received []map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); !strings.HasPrefix(auth, "Splunk ") {
			t.Errorf("missing Splunk auth header, got %q", auth)
		}
		body, _ := io.ReadAll(r.Body)
		dec := json.NewDecoder(strings.NewReader(string(body)))
		for dec.More() {
			var m map[string]any
			_ = dec.Decode(&m)
			received = append(received, m)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSplunkNotifier(ts.URL, "test-token", "portwatch")
	if err := n.Notify([]alert.Event{splunkEvent()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 1 {
		t.Fatalf("expected 1 event, got %d", len(received))
	}
	ev := received[0]["event"].(map[string]any)
	if ev["port"].(float64) != 9200 {
		t.Errorf("expected port 9200, got %v", ev["port"])
	}
}

func TestSplunkNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewSplunkNotifier(ts.URL, "bad-token", "portwatch")
	err := n.Notify([]alert.Event{splunkEvent()})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestSplunkNotifierBadURL(t *testing.T) {
	n := NewSplunkNotifier("://bad-url", "token", "src")
	err := n.Notify([]alert.Event{splunkEvent()})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestSplunkNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewSplunkNotifier(ts.URL, "token", "src")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}
