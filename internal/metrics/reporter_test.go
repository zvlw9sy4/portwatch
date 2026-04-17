package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestReporterPrintContainsFields(t *testing.T) {
	c := metrics.NewCollector()
	c.RecordScan()
	c.RecordAlert()

	var buf bytes.Buffer
	r := metrics.NewReporter(&buf)
	if err := r.Print(c.Snapshot()); err != nil {
		t.Fatalf("Print returned error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"Uptime", "Scans total", "Alerts total", "Last scan", "Last alert"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing field %q\noutput:\n%s", want, out)
		}
	}
}

func TestReporterNeverWhenZero(t *testing.T) {
	s := metrics.Snapshot{
		UptimeSince: time.Now(),
	}
	var buf bytes.Buffer
	r := metrics.NewReporter(&buf)
	if err := r.Print(s); err != nil {
		t.Fatalf("Print returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "never") {
		t.Errorf("expected 'never' for zero times, got:\n%s", buf.String())
	}
}
