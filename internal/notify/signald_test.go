package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func signaldEvent() []alert.Event {
	return []alert.Event{
		{Kind: "opened", Port: scanner.Port{Number: 8080, Proto: "tcp"}},
	}
}

func TestSignaldNotifierSuccess(t *testing.T) {
	var got signaldPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/send" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := NewSignaldNotifier(ts.URL, "+10000000000", []string{"+19999999999"})
	if err := n.Notify(signaldEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Number != "+10000000000" {
		t.Errorf("sender: got %q, want +10000000000", got.Number)
	}
	if len(got.Recipients) != 1 || got.Recipients[0] != "+19999999999" {
		t.Errorf("recipients: got %v", got.Recipients)
	}
	if !strings.Contains(got.Message, "opened") {
		t.Errorf("message missing 'opened': %q", got.Message)
	}
}

func TestSignaldNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewSignaldNotifier(ts.URL, "+10000000000", []string{"+19999999999"})
	if err := n.Notify(signaldEvent()); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestSignaldNotifierBadURL(t *testing.T) {
	n := NewSignaldNotifier("http://127.0.0.1:1", "+1", []string{"+2"})
	if err := n.Notify(signaldEvent()); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestSignaldNotifierNoEventsSkips(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewSignaldNotifier(ts.URL, "+1", []string{"+2"})
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty events")
	}
}
