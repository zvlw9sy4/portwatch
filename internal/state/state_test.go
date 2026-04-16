package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	store := state.NewStore(path)

	snap := state.Snapshot{
		Timestamp: time.Now().Truncate(time.Second),
		Ports: []scanner.Port{
			{Number: 80, Proto: "tcp"},
			{Number: 443, Proto: "tcp"},
		},
	}

	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(got.Ports) != len(snap.Ports) {
		t.Errorf("expected %d ports, got %d", len(snap.Ports), len(got.Ports))
	}
	if !got.Timestamp.Equal(snap.Timestamp) {
		t.Errorf("timestamp mismatch: got %v want %v", got.Timestamp, snap.Timestamp)
	}
}

func TestLoadMissingFile(t *testing.T) {
	store := state.NewStore("/tmp/portwatch_nonexistent_xyz.json")
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if snap.Ports != nil {
		t.Errorf("expected nil ports for empty snapshot")
	}
}

func TestSaveCreatesFile(t *testing.T) {
	path := tempPath(t)
	store := state.NewStore(path)

	if err := store.Save(state.Snapshot{Timestamp: time.Now()}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
