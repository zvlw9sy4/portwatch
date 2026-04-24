// Package notify provides notifier implementations for portwatch alerts.
//
// # Splunk Notifier
//
// SplunkNotifier delivers port-change events to a Splunk HTTP Event Collector
// (HEC) endpoint. Each alert.Event is serialised as a separate HEC JSON event
// so that Splunk indexes them individually with their original timestamps.
//
// Configuration:
//
//	endpoint – full HEC URL, e.g. https://splunk.corp:8088/services/collector
//	token    – HEC token (the value after "Splunk " in the Authorization header)
//	source   – Splunk source field, e.g. "portwatch"
//
// The notifier is a no-op when the event slice is empty, avoiding unnecessary
// HTTP round-trips during quiet polling cycles.
package notify
