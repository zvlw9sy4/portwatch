package alert_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

type recordingNotifier struct {
	events []alert.Event
	failAfter int
	calls     int
}

func (r *recordingNotifier) Notify(e alert.Event) error {
	r.calls++
	if r.failAfter > 0 && r.calls > r.failAfter {
		return errors.New("mock error")
	}
	r.events = append(r.events, e)
	return nil
}

func TestDispatcherFanOut(t *testing.T) {
	a := &recordingNotifier{}
	b := &recordingNotifier{}
	d := alert.NewDispatcher(a, b)

	events := alert.BuildEvents(
		[]scanner.Port{{Number: 3000, Protocol: "tcp"}},
		[]scanner.Port{{Number: 5432, Protocol: "tcp"}},
	)
	d.Dispatch(events)

	if len(a.events) != 2 {
		t.Errorf("notifier a: expected 2 events, got %d", len(a.events))
	}
	if len(b.events) != 2 {
		t.Errorf("notifier b: expected 2 events, got %d", len(b.events))
	}
}

func TestDispatcherContinuesOnError(t *testing.T) {
	failing := &recordingNotifier{failAfter: 0}
	good := &recordingNotifier{}
	d := alert.NewDispatcher(failing, good)

	events := alert.BuildEvents([]scanner.Port{{Number: 80, Protocol: "tcp"}}, nil)
	d.Dispatch(events) // should not panic

	if good.calls != 1 {
		t.Errorf("expected good notifier to receive 1 call, got %d", good.calls)
	}
}

func TestDispatcherAdd(t *testing.T) {
	d := alert.NewDispatcher()
	r := &recordingNotifier{}
	d.Add(r)

	events := alert.BuildEvents(nil, []scanner.Port{{Number: 443, Protocol: "tcp"}})
	d.Dispatch(events)

	if len(r.events) != 1 {
		t.Errorf("expected 1 event after Add, got %d", len(r.events))
	}
}
