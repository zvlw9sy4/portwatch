package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestNewCollectorInitialState(t *testing.T) {
	c := metrics.NewCollector()
	s := c.Snapshot()
	if s.ScansTotal != 0 || s.AlertsTotal != 0 {
		t.Fatalf("expected zero counters, got scans=%d alerts=%d", s.ScansTotal, s.AlertsTotal)
	}
	if s.UptimeSince.IsZero() {
		t.Fatal("UptimeSince should not be zero")
	}
}

func TestRecordScan(t *testing.T) {
	c := metrics.NewCollector()
	before := time.Now()
	c.RecordScan()
	c.RecordScan()
	s := c.Snapshot()
	if s.ScansTotal != 2 {
		t.Fatalf("expected 2 scans, got %d", s.ScansTotal)
	}
	if s.LastScanAt.Before(before) {
		t.Fatal("LastScanAt should be after test start")
	}
}

func TestRecordAlert(t *testing.T) {
	c := metrics.NewCollector()
	c.RecordAlert()
	s := c.Snapshot()
	if s.AlertsTotal != 1 {
		t.Fatalf("expected 1 alert, got %d", s.AlertsTotal)
	}
	if s.LastAlertAt.IsZero() {
		t.Fatal("LastAlertAt should be set")
	}
}

func TestSnapshotIsCopy(t *testing.T) {
	c := metrics.NewCollector()
	s1 := c.Snapshot()
	c.RecordScan()
	s2 := c.Snapshot()
	if s1.ScansTotal == s2.ScansTotal {
		t.Fatal("snapshots should be independent copies")
	}
}
