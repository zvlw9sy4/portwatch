package ratelimit

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllowFirstCallAlwaysPasses(t *testing.T) {
	l := New(time.Minute)
	if !l.Allow("open:tcp:8080") {
		t.Fatal("first call should be allowed")
	}
}

func TestAllowSuppressesWithinCooldown(t *testing.T) {
	base := time.Now()
	l := New(time.Minute)
	l.now = fixedClock(base)

	l.Allow("open:tcp:8080")

	l.now = fixedClock(base.Add(30 * time.Second))
	if l.Allow("open:tcp:8080") {
		t.Fatal("should be suppressed within cooldown")
	}
}

func TestAllowPassesAfterCooldown(t *testing.T) {
	base := time.Now()
	l := New(time.Minute)
	l.now = fixedClock(base)

	l.Allow("open:tcp:8080")

	l.now = fixedClock(base.Add(61 * time.Second))
	if !l.Allow("open:tcp:8080") {
		t.Fatal("should be allowed after cooldown expires")
	}
}

func TestAllowIndependentKeys(t *testing.T) {
	base := time.Now()
	l := New(time.Minute)
	l.now = fixedClock(base)

	l.Allow("open:tcp:8080")

	if !l.Allow("open:tcp:9090") {
		t.Fatal("different key should be allowed")
	}
}

func TestResetAllowsImmediateRetry(t *testing.T) {
	base := time.Now()
	l := New(time.Minute)
	l.now = fixedClock(base)

	l.Allow("open:tcp:8080")
	l.Reset("open:tcp:8080")

	if !l.Allow("open:tcp:8080") {
		t.Fatal("reset key should be allowed immediately")
	}
}

func TestFlushRemovesExpiredEntries(t *testing.T) {
	base := time.Now()
	l := New(time.Minute)
	l.now = fixedClock(base)

	l.Allow("open:tcp:8080")

	l.now = fixedClock(base.Add(2 * time.Minute))
	l.Flush()

	if len(l.last) != 0 {
		t.Fatalf("expected empty map after flush, got %d entries", len(l.last))
	}
}
