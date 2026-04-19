// Package snapshot tracks successive port scan results and produces
// human-readable change summaries.
//
// Usage:
//
//	rec := snapshot.NewRecorder(nil)
//	for {
//		ports := scanner.Scan(...)
//		summary := rec.Update(ports)
//		if len(summary.Opened) > 0 || len(summary.Closed) > 0 {
//			fmt.Println(summary)
//		}
//	}
package snapshot
