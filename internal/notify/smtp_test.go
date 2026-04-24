package notify

import (
	"bufio"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// startFakeSMTPServer starts a minimal fake SMTP server and returns its address.
// It reads the DATA payload into received and signals done when QUIT is seen.
func startFakeSMTPServer(t *testing.T, received *string, done chan struct{}) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		defer ln.Close()
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(5 * time.Second)) //nolint:errcheck
		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		writeLine := func(s string) { rw.WriteString(s + "\r\n"); rw.Flush() } //nolint:errcheck
		readLine := func() string { l, _ := rw.ReadString('\n'); return strings.TrimSpace(l) }
		writeLine("220 fake smtp ready")
		for {
			line := readLine()
			switch {
			case strings.HasPrefix(line, "EHLO"), strings.HasPrefix(line, "HELO"):
				writeLine("250 OK")
			case strings.HasPrefix(line, "AUTH"):
				writeLine("235 OK")
			case strings.HasPrefix(line, "MAIL FROM"):
				writeLine("250 OK")
			case strings.HasPrefix(line, "RCPT TO"):
				writeLine("250 OK")
			case line == "DATA":
				writeLine("354 send data")
				var buf strings.Builder
				for {
					l, _ := rw.ReadString('\n')
					if strings.TrimSpace(l) == "." {
						break
					}
					buf.WriteString(l)
				}
				*received = buf.String()
				writeLine("250 OK")
			case strings.HasPrefix(line, "QUIT"):
				writeLine("221 bye")
				close(done)
				return
			default:
				_ = io.Discard
			}
		}
	}()
	return ln.Addr().String()
}

func smtpEvent() alert.Event {
	return alert.Event{
		Kind: "opened",
		Port: scanner.Port{Number: "8080", Proto: "tcp"},
	}
}

func TestSMTPNotifierSendsMessage(t *testing.T) {
	var received string
	done := make(chan struct{})
	addr := startFakeSMTPServer(t, &received, done)

	host, port := splitHostPort(t, addr)

	n := NewSMTPNotifier(SMTPOptions{
		Host:     host,
		Port:     port,
		Username: "user",
		Password: "pass",
		From:     "alert@portwatch.local",
		To:       []string{"admin@example.com"},
	})

	if err := n.Notify([]alert.Event{smtpEvent()}); err != nil {
		t.Fatalf("Notify: %v", err)
	}

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("server did not finish in time")
	}

	if !strings.Contains(received, "8080") {
		t.Errorf("expected port 8080 in body, got: %s", received)
	}
	if !strings.Contains(received, "opened") {
		t.Errorf("expected 'opened' in body, got: %s", received)
	}
}

func TestSMTPNotifierNoEventsSkipsSend(t *testing.T) {
	n := NewSMTPNotifier(SMTPOptions{
		Host: "127.0.0.1",
		Port: 1,
		From: "a@b.com",
		To:   []string{"c@d.com"},
	})
	// Should not attempt a connection and must not error.
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected no error for empty events, got: %v", err)
	}
}

func TestSMTPNotifierBadHost(t *testing.T) {
	n := NewSMTPNotifier(SMTPOptions{
		Host: "127.0.0.1",
		Port: 1, // nothing listening
		From: "a@b.com",
		To:   []string{"c@d.com"},
	})
	err := n.Notify([]alert.Event{smtpEvent()})
	if err == nil {
		t.Fatal("expected error for unreachable host")
	}
}

// splitHostPort splits an address string into host and integer port for tests.
func splitHostPort(t *testing.T, addr string) (string, int) {
	t.Helper()
	ln, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	return ln.IP.String(), ln.Port
}
