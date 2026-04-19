package notify_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/scanner"
)

func syslogEvent(kind alert.EventKind, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: port, Protocol: "tcp", Address: "0.0.0.0"},
	}
}

func TestSyslogNotifierCreates(t *testing.T) {
	n, err := notify.NewSyslogNotifier("portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer n.Close()

	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSyslogNotifierSendsOpened(t *testing.T) {
	n, err := notify.NewSyslogNotifier("portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer n.Close()

	e := syslogEvent(alert.EventOpened, 8080)
	if err := n.Notify(e); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
}

func TestSyslogNotifierSendsClosed(t *testing.T) {
	n, err := notify.NewSyslogNotifier("portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer n.Close()

	e := syslogEvent(alert.EventClosed, 9090)
	if err := n.Notify(e); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
}
