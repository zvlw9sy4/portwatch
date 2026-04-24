// Package notify provides notifier implementations for portwatch.
//
// # Lark / Feishu Notifier
//
// NewLarkNotifier creates a notifier that posts alerts to a Lark (Feishu)
// incoming-webhook URL.
//
// Usage:
//
//	n := notify.NewLarkNotifier("https://open.feishu.cn/open-apis/bot/v2/hook/<token>")
//
// Messages are sent as plain-text cards using the "text" message type.
// Each event is rendered on its own line:
//
//	[portwatch] OPENED port tcp/8080
//
// If the event list is empty the HTTP call is skipped entirely.
// A non-2xx response is returned as an error.
package notify
