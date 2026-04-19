// Package notify provides notifier implementations for delivering portwatch
// alerts to external services.
//
// Supported backends:
//   - Webhook  — generic HTTP POST
//   - Slack    — Slack incoming webhooks
//   - PagerDuty — PagerDuty Events API v2
//   - Email    — SMTP email delivery
//   - Teams    — Microsoft Teams incoming webhooks
//
// Each notifier implements the alert.Notifier interface and can be registered
// with an alert.Dispatcher to fan alerts out to multiple destinations.
package notify
