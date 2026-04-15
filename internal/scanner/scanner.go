package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents an open port detected on the system.
type Port struct {
	Protocol string
	Address  string
	Port     int
}

// String returns a human-readable representation of the port.
func (p Port) String() string {
	return fmt.Sprintf("%s://%s:%d", p.Protocol, p.Address, p.Port)
}

// Scanner scans for open TCP ports within a given range.
type Scanner struct {
	Host    string
	MinPort int
	MaxPort int
	Timeout time.Duration
}

// NewScanner creates a Scanner with sensible defaults.
func NewScanner(host string, minPort, maxPort int) *Scanner {
	return &Scanner{
		Host:    host,
		MinPort: minPort,
		MaxPort: maxPort,
		Timeout: 500 * time.Millisecond,
	}
}

// Scan performs a TCP connect scan and returns all open ports.
func (s *Scanner) Scan() ([]Port, error) {
	if s.MinPort < 1 || s.MaxPort > 65535 || s.MinPort > s.MaxPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", s.MinPort, s.MaxPort)
	}

	var open []Port
	for port := s.MinPort; port <= s.MaxPort; port++ {
		addr := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err != nil {
			continue
		}
		conn.Close()
		open = append(open, Port{
			Protocol: "tcp",
			Address:  s.Host,
			Port:     port,
		})
	}
	return open, nil
}
