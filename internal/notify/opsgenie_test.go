package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func ogEvent() alert.Event {
	return alert.Event{
		Port:  scanner.Port{Number: 8080, Protocol: "tcp"},
		Kind:  alert.Opened,
	}
}

func TestOpsGenieNotifierSuccess(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("missing Authorization header")
		}
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewOpsGenieNotifier("test-key")
	n.apiURL = ts.URL

	if err := n.Notify(ogEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["message"] == "" {
		t.Error("expected message in payload")
	}
}

func TestOpsGenieNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewOpsGenieNotifier("bad-key")
	n.apiURL = ts.URL

	if err := n.Notify(ogEvent()); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestOpsGenieNotifierBadURL(t *testing.T) {
	n := NewOpsGenieNotifier("key")
	n.apiURL = "http://127.0.0.1:0"

	if err := n.Notify(ogEvent()); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
