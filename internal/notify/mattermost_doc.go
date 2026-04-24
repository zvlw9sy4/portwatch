// Package notify provides notifier implementations for various alerting
// backends used by portwatch.
//
// # Mattermost
//
// MattermostNotifier delivers port-change events to a Mattermost channel via
// an incoming webhook integration.
//
// Set up an incoming webhook in your Mattermost instance under
// Integrations → Incoming Webhooks, then supply the generated URL:
//
//	notifier := notify.NewMattermostNotifier(
//		"https://mattermost.example.com/hooks/xxxxxxxxxxxx",
//		"#security-alerts", // leave empty to use the webhook default
//	)
//
// Each call to Notify batches all events into a single message using
// Mattermost markdown so that opened ports appear with a green indicator
// and closed ports appear with a red indicator.
package notify
