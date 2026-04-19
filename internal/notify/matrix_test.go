package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func matrixEvent() alert.Event {
	return alert.Event{Kind: "opened", Port: "tcp/8080"}
}

func TestMatrixNotifierSuccess(t *testing.T) {
	var gotAuth, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		var m map[string]string
		_ = json.NewDecoder(r.Body).Decode(&m)
		gotBody = m["body"]
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"event_id":"$abc"}`))
	}))
	defer ts.Close()

	n := NewMatrixNotifier(ts.URL, "tok123", "!room:example.com")
	if err := n.Notify(matrixEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Bearer tok123" {
		t.Errorf("expected Bearer token, got %q", gotAuth)
	}
	if gotBody == "" {
		t.Error("expected non-empty message body")
	}
}

func TestMatrixNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewMatrixNotifier(ts.URL, "bad", "!room:example.com")
	if err := n.Notify(matrixEvent()); err == nil {
		t.Error("expected error for non-2xx response")
	}
}

func TestMatrixNotifierBadURL(t *testing.T) {
	n := NewMatrixNotifier("http://127.0.0.1:1", "tok", "!room:example.com")
	if err := n.Notify(matrixEvent()); err == nil {
		t.Error("expected error for unreachable server")
	}
}
