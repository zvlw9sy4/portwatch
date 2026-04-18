// Package ratelimit provides alert throttling to suppress duplicate
// notifications within a configurable cooldown window.
package ratelimit

import (
	"sync"
	"time"
)

// Key identifies a unique alert event (e.g. "open:tcp:8080").
type Key = string

// Limiter suppresses repeated alerts for the same key within a cooldown period.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[Key]time.Time
	now      func() time.Time
}

// New creates a Limiter with the given cooldown duration.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[Key]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the alert for key should be forwarded, and records
// the current time so subsequent calls within the cooldown window return false.
func (l *Limiter) Allow(key Key) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if t, seen := l.last[key]; seen && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded time for key, allowing the next alert through
// regardless of the cooldown window.
func (l *Limiter) Reset(key Key) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// Flush removes all entries whose last-seen time is older than the cooldown,
// preventing unbounded memory growth in long-running daemons.
func (l *Limiter) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	for k, t := range l.last {
		if now.Sub(t) >= l.cooldown {
			delete(l.last, k)
		}
	}
}
