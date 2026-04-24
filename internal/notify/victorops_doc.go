// Package notify provides notifier implementations for portwatch.
//
// # VictorOps (Splunk On-Call) Notifier
//
// VictorOpsNotifier delivers alerts to VictorOps via the REST Endpoint
// integration. Each port-opened event is sent as a CRITICAL message and
// each port-closed event is sent as an INFO message, enabling automatic
// incident creation and resolution in your on-call workflow.
//
// # Usage
//
//	n := notify.NewVictorOpsNotifier(
//	    "https://alert.victorops.com/integrations/generic/12345/alert",
//	    "my-routing-key",
//	)
//
// The routing key is appended to the endpoint URL as a path segment, matching
// the format expected by the VictorOps REST Endpoint integration.
package notify
