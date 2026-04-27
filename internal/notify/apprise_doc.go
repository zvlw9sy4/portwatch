// Package notify provides alert.Notifier implementations for various
// notification back-ends.
//
// # Apprise
//
// AppriseNotifier delivers portwatch alerts to a self-hosted Apprise API
// server (https://github.com/caronc/apprise-api).  Apprise itself acts as
// a fan-out gateway that forwards the message to any of its 80+ supported
// services (Slack, Discord, email, SMS, …) based on the server-side
// configuration.
//
// Usage:
//
//	n := notify.NewAppriseNotifier("http://apprise.internal:8000", "portwatch")
//
// The second argument is the optional Apprise tag that scopes delivery to a
// named group configured on the server.  Pass an empty string to use the
// server default.
//
// The notifier skips the HTTP call entirely when the event list is empty,
// avoiding spurious wake-ups on quiet polling cycles.
package notify
