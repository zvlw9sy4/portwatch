package watch

import (
	"log"
	"sync"

	"github.com/user/portwatch/internal/config"
)

// ReloadableConfig holds the latest config and refreshes it when the file changes.
type ReloadableConfig struct {
	mu   sync.RWMutex
	cfg  config.Config
	path string
}

// NewReloadableConfig loads the initial config and sets up hot-reload via Watcher.
func NewReloadableConfig(path string, interval int) (*ReloadableConfig, *Watcher, error) {
	cfg, err := config.Load(path)
	if err != nil {
		return nil, nil, err
	}
	rc := &ReloadableConfig{cfg: cfg, path: path}
	w := New(path, pollInterval(interval), func() {
		if err := rc.reload(); err != nil {
			log.Printf("[watch] reload failed: %v", err)
		}
	})
	return rc, w, nil
}

// Get returns the current config snapshot.
func (rc *ReloadableConfig) Get() config.Config {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.cfg
}

func (rc *ReloadableConfig) reload() error {
	cfg, err := config.Load(rc.path)
	if err != nil {
		return err
	}
	rc.mu.Lock()
	rc.cfg = cfg
	rc.mu.Unlock()
	log.Printf("[watch] config reloaded from %s", rc.path)
	return nil
}

func pollInterval(seconds int) Duration {
	if seconds <= 0 {
		seconds = 5
	}
	return Duration(seconds) * Second
}
