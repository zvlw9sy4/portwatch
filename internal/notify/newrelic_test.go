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

func newRelicEvent() alert.Event {
	return alert.Event{
		Port: scanner.Port{Number: 9200, Proto: "tcp"},
		Kind: alert.Opened,
	}
}

func TestNewRelicNotifierSuccess(t *testing.T) {
	var gotHeader, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("Api-Key")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewNewRelicNotifier("test-key-123", ts.URL)
	if err := n.Notify([]alert.Event{newRelicEvent()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotHeader != "test-key-123" {
		t.Errorf("expected Api-Key header 'test-key-123', got %q", gotHeader)
	}

	var payload struct {
		Logs []struct {
			Message    string            `json:"message"`
			Attributes map[string]string `json:"attributes"`
		} `json:"logs"`
	}
	if err := json.Unmarshal([]byte(gotBody), &payload); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if len(payload.Logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(payload.Logs))
	}
	if payload.Logs[0].Attributes["port"] != "9200" {
		t.Errorf("expected port attribute '9200', got %q", payload.Logs[0].Attributes["port"])
	}
}

func TestNewRelicNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewNewRelicNotifier("bad-key", ts.URL)
	err := n.Notify([]alert.Event{newRelicEvent()})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestNewRelicNotifierBadURL(t *testing.T) {
	n := NewNewRelicNotifier("key", "http://127.0.0.1:0")
	err := n.Notify([]alert.Event{newRelicEvent()})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestNewRelicNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewNewRelicNotifier("key", ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}
