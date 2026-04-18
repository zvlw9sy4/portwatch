package notify

import (
	"io"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// startFakeSMTP spins up a minimal SMTP server that accepts one message.
func startFakeSMTP(t *testing.T) (addr string, received chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	received = make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		ln.Close()
		if err != nil {
			return
		}
		defer conn.Close()
		tc := textproto.NewConn(conn)
		_ = tc.PrintfLine("220 fake smtp")
		var buf strings.Builder
		for {
			line, err := tc.ReadLine()
			if err == io.EOF {
				break
			}
			upper := strings.ToUpper(line)
			switch {
			case strings.HasPrefix(upper, "EHLO"), strings.HasPrefix(upper, "HELO"):
				_ = tc.PrintfLine("250 ok")
			case strings.HasPrefix(upper, "AUTH"):
				_ = tc.PrintfLine("235 ok")
			case upper == "MAIL FROM:" || strings.HasPrefix(upper, "MAIL"):
				_ = tc.PrintfLine("250 ok")
			case strings.HasPrefix(upper, "RCPT"):
				_ = tc.PrintfLine("250 ok")
			case upper == "DATA":
				_ = tc.PrintfLine("354 go ahead")
			case line == ".":
				_ = tc.PrintfLine("250 ok")
				received <- buf.String()
				return
			case strings.HasPrefix(upper, "QUIT"):
				_ = tc.PrintfLine("221 bye")
				return
			default:
				buf.WriteString(line + "\n")
			}
		}
	}()
	return ln.Addr().String(), received
}

func TestEmailNotifierSendsMessage(t *testing.T) {
	_ = smtp.PlainAuth // ensure package used
	addr, received := startFakeSMTP(t)
	host, port := func() (string, int) {
		h, p, _ := net.SplitHostPort(addr)
		var pi int
		fmt.Sscanf(p, "%d", &pi)
		return h, pi
	}()

	cfg := EmailConfig{
		Host: host, Port: port,
		Username: "u", Password: "p",
		From: "from@example.com",
		To: []string{"to@example.com"},
	}
	n := NewEmailNotifier(cfg)
	ev := alert.Event{
		Kind:     "opened",
		Port:     scanner.Port{Number: 8080},
		Protocol: "tcp",
	}
	if err := n.Notify(ev); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
	msg := <-received
	if !strings.Contains(msg, "8080") {
		t.Errorf("expected port in body, got: %s", msg)
	}
}
