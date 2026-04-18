// Package ratelimit implements a thread-safe cooldown limiter for portwatch
// alert notifications.
//
// A Limiter tracks the last time each alert key was forwarded and suppresses
// subsequent identical alerts until the configured cooldown window elapses.
// This prevents alert storms when a port repeatedly flaps open and closed.
//
// Typical usage:
//
//	limiter := ratelimit.New(5 * time.Minute)
//	if limiter.Allow(key) {
//		notifier.Send(event)
//	}
//
// Call Flush periodically to reclaim memory from expired entries.
package ratelimit
