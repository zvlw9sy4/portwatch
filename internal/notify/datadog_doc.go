// Package notify provides notifier implementations for portwatch.
//
// # Datadog Notifier
//
// DatadogNotifier posts port-change events to the Datadog Events API
// (https://docs.datadoghq.com/api/latest/events/).
//
// Each alert.Event is translated into a Datadog event with:
//   - alert_type "warning" for newly opened ports
//   - alert_type "info" for closed ports
//   - tags: ["source:portwatch", "port:<number/proto>"]
//
// Usage:
//
//	n := notify.NewDatadogNotifier(apiKey, "")
//	err := n.Notify(events)
//
// Pass an empty baseURL to use the default US endpoint
// (https://api.datadoghq.com). For EU customers use
// https://api.datadoghq.eu.
package notify
