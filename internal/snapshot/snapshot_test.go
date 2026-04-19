package snapshot_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func port(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestUpdateNoChanges(t *testing.T) {
	r := snapshot.NewRecorder(fixedClock(time.Unix(0, 0)))
	ports := []scanner.Port{port(80, "tcp")}
	r.Update(ports)
	s := r.Update(ports)
	if len(s.Opened) != 0 || len(s.Closed) != 0 {
		t.Fatalf("expected no changes, got opened=%v closed=%v", s.Opened, s.Closed)
	}
}

func TestUpdateDetectsOpened(t *testing.T) {
	r := snapshot.NewRecorder(nil)
	r.Update([]scanner.Port{port(80, "tcp")})
	s := r.Update([]scanner.Port{port(80, "tcp"), port(443, "tcp")})
	if len(s.Opened) != 1 || s.Opened[0].Number != 443 {
		t.Fatalf("expected port 443 opened, got %v", s.Opened)
	}
	if len(s.Closed) != 0 {
		t.Fatalf("expected no closed ports, got %v", s.Closed)
	}
}

func TestUpdateDetectsClosed(t *testing.T) {
	r := snapshot.NewRecorder(nil)
	r.Update([]scanner.Port{port(80, "tcp"), port(22, "tcp")})
	s := r.Update([]scanner.Port{port(80, "tcp")})
	if len(s.Closed) != 1 || s.Closed[0].Number != 22 {
		t.Fatalf("expected port 22 closed, got %v", s.Closed)
	}
}

func TestSummaryStringNoChanges(t *testing.T) {
	r := snapshot.NewRecorder(fixedClock(time.Unix(0, 0)))
	r.Update([]scanner.Port{port(80, "tcp")})
	s := r.Update([]scanner.Port{port(80, "tcp")})
	if !strings.Contains(s.String(), "no changes") {
		t.Fatalf("expected 'no changes' in output, got: %s", s.String())
	}
}

func TestSummaryStringShowsChanges(t *testing.T) {
	r := snapshot.NewRecorder(fixedClock(time.Unix(0, 0)))
	r.Update(nil)
	s := r.Update([]scanner.Port{port(8080, "tcp")})
	out := s.String()
	if !strings.Contains(out, "opened") {
		t.Fatalf("expected 'opened' in output, got: %s", out)
	}
}
