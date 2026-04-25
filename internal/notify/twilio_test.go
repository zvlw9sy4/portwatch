package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func twilioEvent() []alert.Event {
	return []alert.Event{
		{Kind: alert.Opened, Port: scanner.Port{Number: 22, Protocol: "tcp"}},
	}
}

func TestTwilioNotifierSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if r.FormValue("To") == "" {
			t.Error("missing To field")
		}
		if r.FormValue("Body") == "" {
			t.Error("missing Body field")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := NewTwilioNotifier("ACtest", "token", "+10000000000", "+19999999999")
	n.baseURL = ts.URL

	if err := n.Notify(twilioEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTwilioNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewTwilioNotifier("ACtest", "bad", "+10000000000", "+19999999999")
	n.baseURL = ts.URL

	if err := n.Notify(twilioEvent()); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestTwilioNotifierBadURL(t *testing.T) {
	n := NewTwilioNotifier("ACtest", "token", "+10000000000", "+19999999999")
	n.baseURL = "://invalid"

	if err := n.Notify(twilioEvent()); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestTwilioNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := NewTwilioNotifier("ACtest", "token", "+10000000000", "+19999999999")
	n.baseURL = ts.URL

	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}
