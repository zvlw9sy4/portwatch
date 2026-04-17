package state

import (
	"os"
	"testing"
	"time"
)

func TestIsStaleReturnsFalseForFreshFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	opts := CleanupOptions{MaxAge: time.Hour}
	if IsStale(f.Name(), opts) {
		t.Error("expected fresh file to not be stale")
	}
}

func TestIsStaleReturnsTrueForOldFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	past := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(f.Name(), past, past); err != nil {
		t.Fatal(err)
	}

	opts := CleanupOptions{MaxAge: 24 * time.Hour}
	if !IsStale(f.Name(), opts) {
		t.Error("expected old file to be stale")
	}
}

func TestIsStaleReturnsFalseForMissingFile(t *testing.T) {
	opts := DefaultCleanupOptions()
	if IsStale("/nonexistent/path/state.json", opts) {
		t.Error("expected missing file to return false")
	}
}

func TestRemoveIfStaleDeletesOldFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	past := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(f.Name(), past, past); err != nil {
		t.Fatal(err)
	}

	opts := CleanupOptions{MaxAge: time.Hour}
	removed, err := RemoveIfStale(f.Name(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !removed {
		t.Error("expected file to be removed")
	}
	if _, err := os.Stat(f.Name()); !os.IsNotExist(err) {
		t.Error("expected file to no longer exist")
	}
}

func TestRemoveIfStaleKeepsFreshFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	opts := CleanupOptions{MaxAge: time.Hour}
	removed, err := RemoveIfStale(f.Name(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed {
		t.Error("expected fresh file to be kept")
	}
}
