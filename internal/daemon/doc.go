// Package daemon wires together the scanner, state store, and alert
// dispatcher into a long-running scan loop.
//
// Usage:
//
//	d := daemon.New(cfg, store, dispatch)
//	err := d.Run(ctx) // blocks until ctx is cancelled
//
// On each tick the daemon:
//  1. Scans the configured port range.
//  2. Compares the result against the previously persisted state.
//  3. Dispatches alert events for any opened or closed ports.
//  4. Persists the new state for the next tick.
package daemon
