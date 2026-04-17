// Package metrics provides a lightweight, thread-safe collector for portwatch
// daemon runtime statistics such as scan counts, alert counts, and timestamps.
//
// Usage:
//
//	col := metrics.NewCollector()
//	// ... inside daemon loop:
//	col.RecordScan()
//	col.RecordAlert()
//
//	// Print a summary:
//	rep := metrics.NewReporter(os.Stdout)
//	rep.Print(col.Snapshot())
package metrics
