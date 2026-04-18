// Package watch provides a file-system watcher that reloads config
// when the portwatch config file changes on disk.
package watch

import (
	"log"
	"os"
	"time"
)

// Watcher polls a file for modification and calls onChange when it changes.
type Watcher struct {
	path     string
	interval time.Duration
	onChange func()
	stop     chan struct{}
}

// New creates a Watcher for the given file path.
func New(path string, interval time.Duration, onChange func()) *Watcher {
	return &Watcher{
		path:     path,
		interval: interval,
		onChange: onChange,
		stop:     make(chan struct{}),
	}
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	go w.run()
}

// Stop signals the watcher to exit.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) run() {
	lastMod := w.modTime()
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			if t := w.modTime(); !t.IsZero() && t.After(lastMod) {
				lastMod = t
				log.Printf("[watch] config file changed: %s", w.path)
				w.onChange()
			}
		}
	}
}

func (w *Watcher) modTime() time.Time {
	info, err := os.Stat(w.path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
