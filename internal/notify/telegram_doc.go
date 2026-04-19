// Package notify provides notifier implementations for portwatch alerts.
//
// # Telegram Notifier
//
// TelegramNotifier delivers alerts to a Telegram chat using the Bot API.
// Create a bot via @BotFather, obtain a token, and find your chat ID.
//
// Usage:
//
//	n := notify.NewTelegramNotifier("<bot-token>", "<chat-id>")
//
// The notifier sends a plain-text message in the format:
//
//	[portwatch] opened port tcp/8443
//
// Errors are returned for non-2xx HTTP responses or network failures.
package notify
