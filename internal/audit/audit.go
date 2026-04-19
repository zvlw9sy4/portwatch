// Package audit provides a persistent audit log of port change events.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"` // "opened" or "closed"
	Port      uint16    `json:"port"`
	Protocol  string    `json:"protocol"`
}

// Logger appends audit entries to a newline-delimited JSON file.
type Logger struct {
	mu   sync.Mutex
	path string
}

// NewLogger creates a Logger that writes to path.
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Record appends an entry to the audit log.
func (l *Logger) Record(e Entry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// ReadAll reads all entries from the audit log file.
func ReadAll(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
