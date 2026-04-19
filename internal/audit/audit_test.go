package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/portwatch/internal/audit"
)

func tempLog(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.jsonl")
}

func TestRecordAndReadAll(t *testing.T) {
	path := tempLog(t)
	l := audit.NewLogger(path)

	entries := []audit.Entry{
		{Event: "opened", Port: 8080, Protocol: "tcp"},
		{Event: "closed", Port: 22, Protocol: "tcp"},
	}
	for _, e := range entries {
		if err := l.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	got, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Port != 8080 || got[1].Port != 22 {
		t.Errorf("unexpected ports: %v", got)
	}
}

func TestRecordSetsTimestamp(t *testing.T) {
	path := tempLog(t)
	l := audit.NewLogger(path)

	before := time.Now().UTC()
	if err := l.Record(audit.Entry{Event: "opened", Port: 443, Protocol: "tcp"}); err != nil {
		t.Fatal(err)
	}

	got, _ := audit.ReadAll(path)
	if got[0].Timestamp.Before(before) {
		t.Errorf("timestamp not set correctly: %v", got[0].Timestamp)
	}
}

func TestReadAllMissingFile(t *testing.T) {
	got, err := audit.ReadAll("/nonexistent/audit.jsonl")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if got != nil {
		t.Errorf("expected nil slice, got %v", got)
	}
}

func TestRecordBadPath(t *testing.T) {
	l := audit.NewLogger("/nonexistent/dir/audit.jsonl")
	err := l.Record(audit.Entry{Event: "opened", Port: 80, Protocol: "tcp"})
	if err == nil {
		t.Fatal("expected error writing to bad path")
	}
}

func TestReadAllEmpty(t *testing.T) {
	path := tempLog(t)
	os.WriteFile(path, []byte{}, 0o644)
	got, err := audit.ReadAll(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}
