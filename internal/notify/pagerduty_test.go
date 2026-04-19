package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func pdEvent() alert.Event {
	return alert.Event{
		Kind: "opened",
		Port: scanner.Port{Number: 8080, Protocol: "tcp"},
	}
}

func TestPagerDutyNotifierSuccess(t *testing.T) {
	var received pdPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("test-key")
	n.endpoint = ts.URL

	if err := n.Notify(pdEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.RoutingKey != "test-key" {
		t.Errorf("routing key = %q, want test-key", received.RoutingKey)
	}
	if received.EventAction != "trigger" {
		t.Errorf("event_action = %q, want trigger", received.EventAction)
	}
	if received.Payload.Severity != "warning" {
		t.Errorf("severity = %q, want warning", received.Payload.Severity)
	}
}

func TestPagerDutyNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("bad-key")
	n.endpoint = ts.URL

	if err := n.Notify(pdEvent()); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestPagerDutyNotifierBadURL(t *testing.T) {
	n := NewPagerDutyNotifier("key")
	n.endpoint = "http://127.0.0.1:0"

	if err := n.Notify(pdEvent()); err == nil {
		t.Fatal("expected error for unreachable endpoint")
	}
}
