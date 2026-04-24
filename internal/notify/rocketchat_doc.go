// Package notify provides notifier implementations for delivering
// port-change alerts to external services.
//
// # Rocket.Chat
//
// RocketChatNotifier sends alert events to a Rocket.Chat channel via
// an incoming webhook integration.
//
// Usage:
//
//	n := notify.NewRocketChatNotifier("https://chat.example.com/hooks/TOKEN")
//	if err := n.Notify(events); err != nil {
//		log.Println("rocketchat notify:", err)
//	}
//
// Each event is formatted as a single line in the message body:
//
//	[portwatch] opened port tcp/9200
//
// The notifier skips the HTTP call entirely when the event slice is
// empty, avoiding unnecessary webhook traffic.
package notify
