package watch_test

import (
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cfg*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(f.Name(), []byte(content), 0o644)
	return f.Name()
}

func TestWatcherCallsOnChangeWhenFileModified(t *testing.T) {
	path := writeTempFile(t, "initial")

	var called atomic.Int32
	w := watch.New(path, 20*time.Millisecond, func() { called.Add(1) })
	w.Start()
	defer w.Stop()

	time.Sleep(40 * time.Millisecond)
	// Modify the file.
	_ = os.WriteFile(path, []byte("changed"), 0o644)
	time.Sleep(60 * time.Millisecond)

	if called.Load() == 0 {
		t.Error("expected onChange to be called after file modification")
	}
}

func TestWatcherDoesNotCallOnChangeWhenFileUnchanged(t *testing.T) {
	path := writeTempFile(t, "stable")

	var called atomic.Int32
	w := watch.New(path, 20*time.Millisecond, func() { called.Add(1) })
	w.Start()
	defer w.Stop()

	time.Sleep(80 * time.Millisecond)

	if called.Load() != 0 {
		t.Errorf("expected no onChange calls, got %d", called.Load())
	}
}

func TestWatcherMissingFileDoesNotPanic(t *testing.T) {
	w := watch.New("/nonexistent/path.yaml", 20*time.Millisecond, func() {
		t.Error("onChange should not be called for missing file")
	})
	w.Start()
	time.Sleep(60 * time.Millisecond)
	w.Stop()
}
