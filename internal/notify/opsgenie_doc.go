// Package notify provides notifier implementations for portwatch alerts.
//
// # OpsGenie Notifier
//
// NewOpsGenieNotifier sends portwatch events to OpsGenie as alerts.
//
// Configuration:
//
//	api_key:  <your OpsGenie API key>          # required
//	region:   us                               # "us" (default) or "eu"
//	priority: P3                               # P1–P5, default P3
//	tags:
//	  - portwatch
//
// Each opened-port event creates a new OpsGenie alert. Closed-port events
// automatically resolve (close) the corresponding alert by alias so that
// on-call engineers are not left with stale open alerts.
//
// The notifier uses the GenieKey header for authentication and targets the
// v2 Alerts API endpoint.
package notify
