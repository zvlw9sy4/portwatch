package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func larkEvent() []alert.Event {
	return []alert.Event{
		{
			Kind: alert.Opened,
			Port: scanner.Port{Number: 8080, Protocol: "tcp"},
		},
	}
}

func TestLarkNotifierSuccess(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewLarkNotifier(ts.URL)
	if err := n.Notify(larkEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["msg_type"] != "text" {
		t.Errorf("expected msg_type=text, got %v", received["msg_type"])
	}

	content, _ := received["content"].(map[string]interface{})
	text, _ := content["text"].(string)
	if !strings.Contains(text, "8080") {
		t.Errorf("expected text to contain port 8080, got: %s", text)
	}
}

func TestLarkNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewLarkNotifier(ts.URL)
	if err := n.Notify(larkEvent()); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestLarkNotifierBadURL(t *testing.T) {
	n := NewLarkNotifier("http://127.0.0.1:0/no-server")
	if err := n.Notify(larkEvent()); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestLarkNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewLarkNotifier(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call when events list is empty")
	}
}
