package notify

import (
	"net"
	"strings"
	"testing"

	"github.com/mheads/portwatch/internal/alert"
	"github.com/mheads/portwatch/internal/scanner"
)

func xmppEvent(port int, proto string, opened bool) alert.Event {
	return alert.Event{
		Port:   scanner.Port{Number: port, Protocol: proto},
		Opened: opened,
	}
}

func TestXMPPNotifierNoEventsSkips(t *testing.T) {
	called := false
	n := NewXMPPNotifier("localhost:5222", "bot@example.com", "secret", "admin@example.com")
	n.dial = func(network, addr string) (net.Conn, error) {
		called = true
		return nil, nil
	}
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if called {
		t.Fatal("expected dial not to be called for empty events")
	}
}

func TestXMPPNotifierBadDial(t *testing.T) {
	n := NewXMPPNotifier("bad-host:5222", "bot@example.com", "secret", "admin@example.com")
	n.dial = func(network, addr string) (net.Conn, error) {
		return nil, &net.OpError{Op: "dial", Err: &net.DNSError{Err: "no such host"}}
	}
	events := []alert.Event{xmppEvent(22, "tcp", true)}
	err := n.Notify(events)
	if err == nil {
		t.Fatal("expected error on dial failure")
	}
	if !strings.Contains(err.Error(), "xmpp: dial") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestXMPPBuildBody(t *testing.T) {
	events := []alert.Event{
		xmppEvent(80, "tcp", true),
		xmppEvent(443, "tcp", false),
	}
	body := buildXMPPBody(events)
	if !strings.Contains(body, "OPENED") {
		t.Error("expected OPENED in body")
	}
	if !strings.Contains(body, "CLOSED") {
		t.Error("expected CLOSED in body")
	}
	if !strings.Contains(body, "portwatch alert") {
		t.Error("expected header in body")
	}
}

func TestDomainOf(t *testing.T) {
	if got := domainOf("user@example.com"); got != "example.com" {
		t.Errorf("domainOf: got %q, want %q", got, "example.com")
	}
	if got := domainOf("nodomain"); got != "nodomain" {
		t.Errorf("domainOf no-at: got %q, want %q", got, "nodomain")
	}
}
