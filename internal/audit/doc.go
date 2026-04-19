// Package audit provides append-only audit logging for portwatch.
//
// Entries are written as newline-delimited JSON (NDJSON) to a file on disk.
// Each entry captures the timestamp, event type (opened/closed), port number,
// and protocol. The Logger is safe for concurrent use.
//
// Example usage:
//
//	l := audit.NewLogger("/var/log/portwatch/audit.jsonl")
//	l.Record(audit.Entry{
//		Event:    "opened",
//		Port:     8080,
//		Protocol: "tcp",
//	})
package audit
