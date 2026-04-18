package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestWebhookNotifierSuccess(t *testing.T) {
	var received notify.WebhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewWebhookNotifier(ts.URL)
	payload := notify.WebhookPayload{
		Timestamp: "2024-01-01T00:00:00Z",
		Event:     "opened",
		Port:      8080,
		Protocol:  "tcp",
	}
	if err := n.Notify(payload); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.Event != "opened" {
		t.Errorf("expected event opened, got %s", received.Event)
	}
}

func TestWebhookNotifierNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewWebhookNotifier(ts.URL)
	err := n.Notify(notify.WebhookPayload{Event: "closed", Port: 22, Protocol: "tcp"})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestWebhookNotifierBadURL(t *testing.T) {
	n := notify.NewWebhookNotifier("http://127.0.0.1:0")
	err := n.Notify(notify.WebhookPayload{Event: "opened", Port: 9999, Protocol: "tcp"})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
