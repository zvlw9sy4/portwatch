// Package snapshot provides periodic port snapshot diffing with summary reporting.
package snapshot

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Summary holds a diff summary between two snapshots.
type Summary struct {
	Timestamp time.Time
	Opened    []scanner.Port
	Closed    []scanner.Port
}

// String returns a human-readable summary.
func (s Summary) String() string {
	if len(s.Opened) == 0 && len(s.Closed) == 0 {
		return fmt.Sprintf("[%s] no changes detected", s.Timestamp.Format(time.RFC3339))
	}
	result := fmt.Sprintf("[%s] changes detected:\n", s.Timestamp.Format(time.RFC3339))
	for _, p := range s.Opened {
		result += fmt.Sprintf("  + opened %s\n", p)
	}
	for _, p := range s.Closed {
		result += fmt.Sprintf("  - closed %s\n", p)
	}
	return result
}

// Recorder keeps the previous snapshot and produces summaries on each update.
type Recorder struct {
	prev []scanner.Port
	now  func() time.Time
}

// NewRecorder returns a Recorder. If clockFn is nil, time.Now is used.
func NewRecorder(clockFn func() time.Time) *Recorder {
	if clockFn == nil {
		clockFn = time.Now
	}
	return &Recorder{now: clockFn}
}

// Update compares current ports against the previous snapshot and returns a Summary.
func (r *Recorder) Update(current []scanner.Port) Summary {
	opened, closed := diff(r.prev, current)
	s := Summary{
		Timestamp: r.now(),
		Opened:    opened,
		Closed:    closed,
	}
	r.prev = current
	return s
}

func diff(prev, curr []scanner.Port) (opened, closed []scanner.Port) {
	prevSet := toSet(prev)
	currSet := toSet(curr)
	for _, p := range curr {
		if !prevSet[p.String()] {
			opened = append(opened, p)
		}
	}
	for _, p := range prev {
		if !currSet[p.String()] {
			closed = append(closed, p)
		}
	}
	return
}

func toSet(ports []scanner.Port) map[string]bool {
	m := make(map[string]bool, len(ports))
	for _, p := range ports {
		m[p.String()] = true
	}
	return m
}
