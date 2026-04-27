package notify

import (
	"context"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

type fakeKafkaProducer struct {
	published [][]byte
	errOnce   bool
	closed    bool
}

func (f *fakeKafkaProducer) Publish(_ context.Context, _ string, payload []byte) error {
	if f.errOnce {
		f.errOnce = false
		return errors.New("broker unavailable")
	}
	f.published = append(f.published, payload)
	return nil
}

func (f *fakeKafkaProducer) Close() error {
	f.closed = true
	return nil
}

func kafkaEvent(kind string, port int) alert.Event {
	return alert.Event{Kind: kind, Port: scanner.Port{Number: port, Protocol: "tcp"}}
}

func TestKafkaNotifierSuccess(t *testing.T) {
	prod := &fakeKafkaProducer{}
	n := NewKafkaNotifier(prod, "portwatch")

	events := []alert.Event{kafkaEvent("opened", 8080), kafkaEvent("closed", 9090)}
	if err := n.Notify(context.Background(), events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prod.published) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(prod.published))
	}
}

func TestKafkaNotifierNoEventsSkips(t *testing.T) {
	prod := &fakeKafkaProducer{}
	n := NewKafkaNotifier(prod, "portwatch")

	if err := n.Notify(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prod.published) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(prod.published))
	}
}

func TestKafkaNotifierPublishError(t *testing.T) {
	prod := &fakeKafkaProducer{errOnce: true}
	n := NewKafkaNotifier(prod, "portwatch")

	events := []alert.Event{kafkaEvent("opened", 8080)}
	err := n.Notify(context.Background(), events)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestKafkaNotifierPayloadContainsFields(t *testing.T) {
	prod := &fakeKafkaProducer{}
	n := NewKafkaNotifier(prod, "portwatch")

	if err := n.Notify(context.Background(), []alert.Event{kafkaEvent("opened", 443)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw := string(prod.published[0])
	for _, want := range []string{"\"port\":443", "\"kind\":\"opened\"", "\"protocol\":\"tcp\"", "timestamp"} {
		if !containsStr(raw, want) {
			t.Errorf("payload missing %q; got: %s", want, raw)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}())
}
