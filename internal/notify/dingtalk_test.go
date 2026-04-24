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

func dingtalkEvent() alert.Event {
	return alert.Event{
		Kind: "opened",
		Port: scanner.Port{Number: 8080, Proto: "tcp"},
	}
}

func TestDingTalkNotifierSuccess(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewDingTalkNotifier(ts.URL)
	if err := n.Notify([]alert.Event{dingtalkEvent()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["msgtype"] != "text" {
		t.Errorf("expected msgtype=text, got %v", received["msgtype"])
	}
	textMap, ok := received["text"].(map[string]interface{})
	if !ok {
		t.Fatal("missing text field")
	}
	if !strings.Contains(textMap["content"].(string), "8080") {
		t.Errorf("expected port 8080 in content, got %v", textMap["content"])
	}
}

func TestDingTalkNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewDingTalkNotifier(ts.URL)
	if err := n.Notify([]alert.Event{dingtalkEvent()}); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestDingTalkNotifierBadURL(t *testing.T) {
	n := NewDingTalkNotifier("http://127.0.0.1:0/no-server")
	if err := n.Notify([]alert.Event{dingtalkEvent()}); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestDingTalkNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewDingTalkNotifier(ts.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty event list")
	}
}
