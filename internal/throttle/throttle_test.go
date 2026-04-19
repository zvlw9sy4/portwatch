package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

type fakeClock struct {
	now time.Time
}

func (f *fakeClock) Now() time.Time { return f.now }
func (f *fakeClock) Advance(d time.Duration) { f.now = f.now.Add(d) }

func TestReadyFirstCallAlwaysPasses(t *testing.T) {
	fc := &fakeClock{now: time.Unix(1000, 0)}
	th := throttle.NewWithClock(5*time.Second, fc.Now)
	if !th.Ready() {
		t.Fatal("expected first call to Ready to return true")
	}
}

func TestReadySuppressesWithinInterval(t *testing.T) {
	fc := &fakeClock{now: time.Unix(1000, 0)}
	th := throttle.NewWithClock(5*time.Second, fc.Now)
	th.Ready() // consume first tick
	fc.Advance(3 * time.Second)
	if th.Ready() {
		t.Fatal("expected Ready to return false within interval")
	}
}

func TestReadyPassesAfterInterval(t *testing.T) {
	fc := &fakeClock{now: time.Unix(1000, 0)}
	th := throttle.NewWithClock(5*time.Second, fc.Now)
	th.Ready()
	fc.Advance(5 * time.Second)
	if !th.Ready() {
		t.Fatal("expected Ready to return true after interval elapsed")
	}
}

func TestResetAllowsImmediateReady(t *testing.T) {
	fc := &fakeClock{now: time.Unix(1000, 0)}
	th := throttle.NewWithClock(5*time.Second, fc.Now)
	th.Ready()
	th.Reset()
	if !th.Ready() {
		t.Fatal("expected Ready to return true after Reset")
	}
}

func TestLastRunUpdated(t *testing.T) {
	fc := &fakeClock{now: time.Unix(1000, 0)}
	th := throttle.NewWithClock(5*time.Second, fc.Now)
	if !th.LastRun().IsZero() {
		t.Fatal("expected LastRun to be zero before first tick")
	}
	th.Ready()
	if th.LastRun() != fc.Now() {
		t.Fatalf("expected LastRun %v, got %v", fc.Now(), th.LastRun())
	}
}
