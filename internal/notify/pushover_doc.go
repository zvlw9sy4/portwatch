// Package notify provides notifier implementations for portwatch alerts.
//
// # Pushover Notifier
//
// PushoverNotifier delivers port-change alerts via the Pushover service
// (https://pushover.net). Each event is sent as an individual push
// notification to the configured user or group key.
//
// # Configuration
//
//	- Token: Pushover application API token (required)
//	- User:  Pushover user or group key   (required)
//
// # Example
//
//	n := notify.NewPushoverNotifier("your-app-token", "your-user-key")
//	err := n.Notify(events)
package notify
