package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func appriseEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: port, Protocol: proto},
	}
}

func TestAppriseNotifierSuccess(t *testing.T) {
	var received apprisePayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewAppriseNotifier(srv.URL, "portwatch")
	events := []alert.Event{appriseEvent("opened", "tcp", 8080)}

	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Title == "" {
		t.Error("expected non-empty title")
	}
	if received.Tag != "portwatch" {
		t.Errorf("tag: got %q, want %q", received.Tag, "portwatch")
	}
	if received.Body == "" {
		t.Error("expected non-empty body")
	}
}

func TestAppriseNotifierNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := NewAppriseNotifier(srv.URL, "")
	events := []alert.Event{appriseEvent("closed", "tcp", 22)}

	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestAppriseNotifierBadURL(t *testing.T) {
	n := NewAppriseNotifier("http://127.0.0.1:0", "")
	events := []alert.Event{appriseEvent("opened", "tcp", 443)}

	if err := n.Notify(events); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}

func TestAppriseNotifierNoEventsSkips(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewAppriseNotifier(srv.URL, "")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected HTTP call to be skipped for empty event list")
	}
}
