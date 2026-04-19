// Package throttle implements a simple interval-based throttler used by the
// portwatch daemon to prevent scans from running more frequently than the
// configured interval, even when the event loop ticks faster than expected.
//
// Usage:
//
//	th := throttle.New(30 * time.Second)
//	for {
//		time.Sleep(time.Second)
//		if th.Ready() {
//			runScan()
//		}
//	}
package throttle
