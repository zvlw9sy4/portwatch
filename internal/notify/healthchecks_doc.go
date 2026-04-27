// Package notify provides notifier implementations for portwatch.
//
// # HealthChecks Notifier
//
// HealthChecksNotifier integrates with healthchecks.io (or any compatible
// self-hosted alternative such as Cabot or a custom ping endpoint).
//
// Each port-change event triggers a POST to the configured ping URL so that
// the monitoring service can record the failure signal and fire its own
// downstream alerts (email, SMS, etc.).
//
// # Configuration
//
//	notifier: healthchecks
//	healthchecks_url: https://hc-ping.com/<your-check-uuid>/fail
//
// The URL should point to the "/fail" variant of your check so that every
// unexpected port change is recorded as a failure.  When portwatch runs a
// clean scan with no changes you may separately POST to the base UUID URL to
// signal a successful heartbeat.
package notify
