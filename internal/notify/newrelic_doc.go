// Package notify provides notifier implementations for portwatch.
//
// # New Relic Notifier
//
// NewRelicNotifier delivers port-change events to the New Relic Logs API
// (https://docs.newrelic.com/docs/logs/log-api/introduction-log-api/).
//
// Each portwatch event is translated into a structured log entry with the
// following attributes:
//
//   - port     – the port number as a string
//   - protocol – "tcp" or "udp"
//   - kind     – "opened" or "closed"
//   - source   – always "portwatch"
//
// # Configuration
//
//	notifier := notify.NewNewRelicNotifier(
//	    os.Getenv("NEW_RELIC_LICENSE_KEY"),
//	    "", // empty string uses the default US endpoint
//	)
//
// EU customers should pass "https://log-api.eu.newrelic.com/log/v1" as the
// endpoint parameter.
package notify
