package notify

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

type stubNotifier struct {
	called bool
	err    error
}

func (s *stubNotifier) Notify(alert.Event) error {
	s.called = true
	return s.err
}

func testEvent() alert.Event {
	return alert.Event{Port: scanner.Port{Number: 443, Protocol: "tcp"}, Kind: alert.Opened}
}

func TestMultiNotifierCallsAll(t *testing.T) {
	a, b := &stubNotifier{}, &stubNotifier{}
	m := NewMultiNotifier(a, b)
	if err := m.Notify(testEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.called || !b.called {
		t.Error("expected both notifiers to be called")
	}
}

func TestMultiNotifierCollectsErrors(t *testing.T) {
	a := &stubNotifier{err: errors.New("fail-a")}
	b := &stubNotifier{err: errors.New("fail-b")}
	m := NewMultiNotifier(a, b)
	err := m.Notify(testEvent())
	if err == nil {
		t.Fatal("expected combined error")
	}
	if !b.called {
		t.Error("second notifier should still be called after first fails")
	}
}

func TestMultiNotifierAdd(t *testing.T) {
	m := NewMultiNotifier()
	s := &stubNotifier{}
	m.Add(s)
	m.Notify(testEvent())
	if !s.called {
		t.Error("added notifier was not called")
	}
}
