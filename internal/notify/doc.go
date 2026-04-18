// Package notify provides outbound notification backends for portwatch alerts.
//
// Currently supported backends:
//
//   - WebhookNotifier: HTTP POST JSON payloads to a configured URL.
//
// Each notifier is independent and can be wired into the alert dispatcher
// via the alert.Notifier interface.
package notify
