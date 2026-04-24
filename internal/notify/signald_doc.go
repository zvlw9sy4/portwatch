// Package notify provides notifier implementations for portwatch alerts.
//
// # Signald / signal-cli-rest-api notifier
//
// SignaldNotifier delivers alerts via the signal-cli-rest-api HTTP gateway,
// which wraps the Signal messenger protocol. You need a running instance of
// signal-cli-rest-api (https://github.com/bbernhard/signal-cli-rest-api) and
// a registered Signal number.
//
// Usage:
//
//	n := notify.NewSignaldNotifier(
//		"http://localhost:8080",  // base URL of signal-cli-rest-api
//		"+15550001111",           // sender number (must be registered)
//		[]string{"+15550002222"}, // recipient numbers
//	)
//
// All events in a single tick are batched into one message to avoid flooding
// the Signal account with individual notifications.
package notify
