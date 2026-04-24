// Package notify provides notifier implementations for portwatch alerts.
//
// # HipChat Notifier
//
// HipChatNotifier delivers port-change events to a HipChat room using the
// HipChat v2 REST API (compatible with self-hosted HipChat Server and
// third-party forks such as Stride).
//
// Usage:
//
//	n := notify.NewHipChatNotifier(
//		"room-id-or-name",
//		"personal-access-token",
//		"https://api.hipchat.com",
//	)
//
// Each alert.Event is sent as a separate room notification. Opened-port events
// use a green card colour; closed-port events use yellow. The notify flag is
// always set to true so that room members receive a push notification.
//
// An empty event slice is a no-op — no HTTP request is made.
package notify
