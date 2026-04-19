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

func teamsEvent() alert.Event {
	return alert.Event{
		Kind: "opened",
		Port: scanner.Port{Number: 8080, Proto: "tcp"},
	}
}

func TestTeamsNotifierSuccess(t *testing.T) {
	var gotBody map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL)
	if err := n.Notify(teamsEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotBody["text"], "8080") {
		t.Errorf("expected port in message, got: %s", gotBody["text"])
	}
}

func TestTeamsNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL)
	if err := n.Notify(teamsEvent()); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestTeamsNotifierBadURL(t *testing.T) {
	n := NewTeamsNotifier("http://127.0.0.1:0")
	if err := n.Notify(teamsEvent()); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
