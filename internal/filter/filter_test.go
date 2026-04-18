package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func ports(nums ...int) []scanner.Port {
	var out []scanner.Port
	for _, n := range nums {
		out = append(out, scanner.Port{Number: n, Protocol: "tcp"})
	}
	return out
}

func TestApplyNoRules(t *testing.T) {
	f := filter.New(nil)
	in := ports(80, 443, 8080)
	got := f.Apply(in)
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}

func TestApplyFiltersMatchingPort(t *testing.T) {
	f := filter.New([]filter.Rule{{Port: 80, Protocol: "tcp"}})
	got := f.Apply(ports(80, 443))
	if len(got) != 1 || got[0].Number != 443 {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestApplyProtocolMismatchKept(t *testing.T) {
	f := filter.New([]filter.Rule{{Port: 80, Protocol: "udp"}})
	got := f.Apply(ports(80))
	if len(got) != 1 {
		t.Fatalf("port should be kept when protocol differs")
	}
}

func TestApplyEmptyProtocolMatchesAny(t *testing.T) {
	f := filter.New([]filter.Rule{{Port: 443, Protocol: ""}})
	got := f.Apply(ports(443, 8443))
	if len(got) != 1 || got[0].Number != 8443 {
		t.Fatalf("unexpected result: %v", got)
	}
}
