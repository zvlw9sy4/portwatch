package metrics

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Reporter writes human-readable metric summaries to an io.Writer.
type Reporter struct {
	out io.Writer
}

// NewReporter creates a Reporter that writes to out.
func NewReporter(out io.Writer) *Reporter {
	return &Reporter{out: out}
}

// Print writes a formatted summary of the given Snapshot.
func (r *Reporter) Print(s Snapshot) error {
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	uptime := time.Since(s.UptimeSince).Truncate(time.Second)

	lines := []struct{ k, v string }{
		{"Uptime", uptime.String()},
		{"Scans total", fmt.Sprintf("%d", s.ScansTotal)},
		{"Alerts total", fmt.Sprintf("%d", s.AlertsTotal)},
		{"Last scan", formatTime(s.LastScanAt)},
		{"Last alert", formatTime(s.LastAlertAt)},
	}

	for _, l := range lines {
		if _, err := fmt.Fprintf(w, "%s\t%s\n", l.k, l.v); err != nil {
			return err
		}
	}
	return w.Flush()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	return t.Format(time.RFC3339)
}
