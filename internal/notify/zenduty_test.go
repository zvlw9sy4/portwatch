package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func zendutyEvent() alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Number: 8080, Proto: "tcp"},
	}
}

func TestZendutyNotifierSuccess(t *testing.T) {
	var received zendutyPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	n := NewZendutyNotifier("test-key")
	n.client = &http.Client{Transport: rewriteTransport(srv.URL)}

	if err := n.Notify([]alert.Event{zendutyEvent()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.AlertType != "critical" {
		t.Errorf("expected alert_type=critical, got %q", received.AlertType)
	}
	if received.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestZendutyNotifierNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	n := NewZendutyNotifier("bad-key")
	n.client = &http.Client{Transport: rewriteTransport(srv.URL)}

	err := n.Notify([]alert.Event{zendutyEvent()})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestZendutyNotifierBadURL(t *testing.T) {
	n := NewZendutyNotifier("key")
	n.client = &http.Client{Transport: rewriteTransport("http://127.0.0.1:1")}

	err := n.Notify([]alert.Event{zendutyEvent()})
	if err == nil {
		t.Fatal("expected error for unreachable host")
	}
}

func TestZendutyNotifierNoEventsSkips(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	n := NewZendutyNotifier("test-key")
	n.client = &http.Client{Transport: rewriteTransport(srv.URL)}

	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}
