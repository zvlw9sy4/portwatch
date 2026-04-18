// Package watch implements config file hot-reloading for portwatch.
//
// A Watcher polls a file path at a configurable interval and invokes
// a callback whenever the file's modification time advances. The
// ReloadableConfig type wraps a config.Config and uses a Watcher to
// transparently refresh the in-memory configuration without restarting
// the daemon process.
//
// Typical usage:
//
//	rc, w, err := watch.NewReloadableConfig("config.yaml", 5)
//	if err != nil { log.Fatal(err) }
//	w.Start()
//	defer w.Stop()
//	// rc.Get() always returns the latest config.
package watch
