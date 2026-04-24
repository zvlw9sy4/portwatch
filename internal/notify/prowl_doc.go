// Package notify provides notifier implementations for portwatch.
//
// # Prowl Notifier
//
// ProwlNotifier delivers push notifications to iOS devices using the
// Prowl HTTP API (https://www.prowlapp.com).
//
// Usage:
//
//	n := notify.NewProwlNotifier(apiKey, appName, priority)
//
// Parameters:
//   - apiKey:   your Prowl API key (required)
//   - appName:  application name shown in the notification (default: "portwatch")
//   - priority: integer from -2 (very low) to 2 (emergency); clamped automatically
//
// Each call to Notify formats all events into a single description field
// and posts them to the Prowl API endpoint in one request.
package notify
