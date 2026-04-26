package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func customEventPort() scanner.Port {
	return scanner.Port{Number: 9090, Protocol: "tcp"}
}

func TestCustomEventNotifierSuccess(t *testing.T) {
	var received customEventPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewCustomEventNotifier(ts.URL, "")
	events := []alert.Event{{Kind: alert.Opened, Port: customEventPort()}}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Port != 9090 {
		t.Errorf("expected port 9090, got %d", received.Port)
	}
	if received.EventType != string(alert.Opened) {
		t.Errorf("expected event_type %q, got %q", alert.Opened, received.EventType)
	}
}

func TestCustomEventNotifierSetsAuthHeader(t *testing.T) {
	var authHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := NewCustomEventNotifier(ts.URL, "secret-key")
	events := []alert.Event{{Kind: alert.Closed, Port: customEventPort()}}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if authHeader != "Bearer secret-key" {
		t.Errorf("expected Bearer auth header, got %q", authHeader)
	}
}

func TestCustomEventNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewCustomEventNotifier(ts.URL, "")
	events := []alert.Event{{Kind: alert.Opened, Port: customEventPort()}}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestCustomEventNotifierNoEventsSkips(t *testing.T) {
	n := NewCustomEventNotifier("http://127.0.0.1:0", "")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected no error for empty events, got %v", err)
	}
}

func TestCustomEventNotifierBadURL(t *testing.T) {
	n := NewCustomEventNotifier("://bad-url", "")
	events := []alert.Event{{Kind: alert.Opened, Port: customEventPort()}}
	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
