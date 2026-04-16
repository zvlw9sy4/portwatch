// Package state provides persistence for port scan snapshots.
//
// A Snapshot captures the set of open ports observed at a point in time.
// The Store type serialises snapshots to a JSON file on disk so that the
// daemon can compare the current scan against the previous one across
// restarts, enabling accurate opened/closed port detection even after the
// process is restarted.
//
// Typical usage:
//
//	store := state.NewStore("/var/lib/portwatch/state.json")
//	prev, _ := store.Load()
//	// ... run scan ...
//	next := state.Snapshot{Timestamp: time.Now(), Ports: ports}
//	store.Save(next)
package state
