// Package notify provides notifier implementations for portwatch alerts.
//
// # Zenduty Notifier
//
// ZendutyNotifier sends port-change events to a Zenduty integration using
// the Zenduty Events API (https://docs.zenduty.com/docs/api).
//
// # Usage
//
//	n := notify.NewZendutyNotifier("<integration-key>")
//	err := n.Notify(events)
//
// # Alert types
//
// Events of type alert.Opened are sent with alert_type "critical".
// Events of type alert.Closed are sent with alert_type "info".
//
// If the event list is empty, no HTTP request is made.
package notify
