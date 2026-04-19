// Package throttle provides a scan throttler that enforces a minimum
// interval between successive port scans to avoid CPU spikes.
package throttle

import (
	"sync"
	"time"
)

// Clock abstracts time for testing.
type Clock func() time.Time

// Throttler enforces a minimum interval between calls to Ready.
type Throttler struct {
	mu       sync.Mutex
	interval time.Duration
	lastRun  time.Time
	clock    Clock
}

// New returns a Throttler that allows a tick no more than once per interval.
func New(interval time.Duration) *Throttler {
	return NewWithClock(interval, time.Now)
}

// NewWithClock returns a Throttler using the provided clock (useful in tests).
func NewWithClock(interval time.Duration, clock Clock) *Throttler {
	return &Throttler{
		interval: interval,
		clock:    clock,
	}
}

// Ready returns true if enough time has elapsed since the last scan.
// When it returns true it also resets the internal timer.
func (t *Throttler) Ready() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	if t.lastRun.IsZero() || now.Sub(t.lastRun) >= t.interval {
		t.lastRun = now
		return true
	}
	return false
}

// LastRun returns the timestamp of the most recent accepted tick.
func (t *Throttler) LastRun() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastRun
}

// Reset clears the last-run timestamp so the next call to Ready always passes.
func (t *Throttler) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastRun = time.Time{}
}
