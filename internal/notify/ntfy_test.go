package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func ntfyEvent(kind alert.EventKind) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: 8080, Protocol: "tcp"},
	}
}

func TestNtfyNotifierSuccess(t *testing.T) {
	var gotTitle, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTitle = r.Header.Get("Title")
		buf := make([]byte, 256)
		n, _ := r.Body.Read(buf)
		gotBody = string(buf[:n])
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	if err := n.Notify(ntfyEvent(alert.Opened)); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotTitle != "Port opened" {
		t.Errorf("expected title 'Port opened', got %q", gotTitle)
	}
	if gotBody == "" {
		t.Error("expected non-empty body")
	}
}

func TestNtfyNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	if err := n.Notify(ntfyEvent(alert.Closed)); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestNtfyNotifierBadURL(t *testing.T) {
	n := NewNtfyNotifier("http://127.0.0.1:0", "portwatch")
	if err := n.Notify(ntfyEvent(alert.Opened)); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}

func TestNtfyNotifierTopicInURL(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "alerts")
	_ = n.Notify(ntfyEvent(alert.Opened))
	if gotPath != "/alerts" {
		t.Errorf("expected path '/alerts', got %q", gotPath)
	}
}
