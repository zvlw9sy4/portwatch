package state

import (
	"os"
	"time"
)

// CleanupOptions configures old state file removal behavior.
type CleanupOptions struct {
	// MaxAge is the maximum age of a state file before it is considered stale.
	MaxAge time.Duration
}

// DefaultCleanupOptions returns sensible defaults for cleanup.
func DefaultCleanupOptions() CleanupOptions {
	return CleanupOptions{
		MaxAge: 7 * 24 * time.Hour,
	}
}

// IsStale reports whether the state file at path is older than opts.MaxAge.
// Returns false if the file does not exist or cannot be stat'd.
func IsStale(path string, opts CleanupOptions) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) > opts.MaxAge
}

// RemoveIfStale deletes the state file at path if it is stale.
// Returns true if the file was removed, false otherwise.
func RemoveIfStale(path string, opts CleanupOptions) (bool, error) {
	if !IsStale(path, opts) {
		return false, nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
