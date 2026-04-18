// Package notify provides pluggable notifier backends for portwatch alerts.
//
// Each notifier implements a Notify(alert.Event) error method and can be
// registered with the alert.Dispatcher to receive port-change events.
//
// Available notifiers:
//   - WebhookNotifier – HTTP POST to a configurable endpoint.
//   - EmailNotifier   – SMTP email via PlainAuth.
package notify
