// Package notify provides notifier implementations for portwatch alerts.
//
// SyslogNotifier
//
// SyslogNotifier forwards port-change events to the local syslog daemon using
// Go's standard log/syslog package. Opened-port events are sent at the ALERT
// priority; closed-port events are sent at INFO priority.
//
// Usage:
//
//	n, err := notify.NewSyslogNotifier("portwatch")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer n.Close()
//
//	// Register with a MultiNotifier or use directly.
//	multi.Add(n)
//
// The notifier is not available on Windows; builds on unsupported platforms
// will fail at compile time via the syslog package constraints.
package notify
