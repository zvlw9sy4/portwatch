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
// Alert keys are arbitrary strings that identify a unique alert condition,
// such as a combination of host, port, and state (e.g. "host:port:open").
// The caller is responsible for constructing keys consistently.
//
// Call Flush periodically to reclaim memory from expired entries. A
// reasonable interval is the same duration as the cooldown window itself.
package ratelimit
