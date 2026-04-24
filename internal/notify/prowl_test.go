package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func prowlEvent() []alert.Event {
	return []alert.Event{
		{Kind: alert.Opened, Port: scanner.Port{Number: 9090, Protocol: "tcp"}},
	}
}

func TestProwlNotifierSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if r.FormValue("apikey") == "" {
			t.Error("expected apikey to be set")
		}
		if r.FormValue("application") != "portwatch" {
			t.Errorf("unexpected application: %s", r.FormValue("application"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewProwlNotifier("test-api-key", "portwatch", 0)
	n.httpClient = ts.Client()
	// Override URL by using a custom transport that rewrites host.
	n.httpClient = &http.Client{
		Transport: rewriteTransport(ts.URL),
	}

	if err := n.Notify(prowlEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProwlNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewProwlNotifier("bad-key", "", 0)
	n.httpClient = &http.Client{Transport: rewriteTransport(ts.URL)}

	if err := n.Notify(prowlEvent()); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestProwlNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewProwlNotifier("key", "app", 1)
	n.httpClient = &http.Client{Transport: rewriteTransport(ts.URL)}

	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}

func TestProwlNotifierBadURL(t *testing.T) {
	n := NewProwlNotifier("key", "app", 0)
	n.httpClient = &http.Client{}
	// Point at an invalid address to force a transport error.
	// We swap the const at runtime via the struct field approach instead.
	// Use a server that immediately closes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	ts.Close() // closed immediately

	n.httpClient = &http.Client{Transport: rewriteTransport(ts.URL)}
	if err := n.Notify(prowlEvent()); err == nil {
		t.Fatal("expected error for closed server")
	}
}
