package scanner

import (
	"net"
	"strconv"
	"testing"
	"time"
)

// startTestServer opens a TCP listener on a random port and returns its port number.
func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScanFindsOpenPort(t *testing.T) {
	port, stop := startTestServer(t)
	defer stop()

	s := NewScanner("127.0.0.1", port, port)
	s.Timeout = 200 * time.Millisecond

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(ports))
	}
	if ports[0].Port != port {
		t.Errorf("expected port %d, got %d", port, ports[0].Port)
	}
	if ports[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", ports[0].Protocol)
	}
}

func TestScanInvalidRange(t *testing.T) {
	s := NewScanner("127.0.0.1", 9000, 8000)
	_, err := s.Scan()
	if err == nil {
		t.Fatal("expected error for invalid port range, got nil")
	}
}

func TestPortString(t *testing.T) {
	p := Port{Protocol: "tcp", Address: "127.0.0.1", Port: 8080}
	expected := "tcp://127.0.0.1:8080"
	if p.String() != expected {
		t.Errorf("expected %q, got %q", expected, p.String())
	}
}
