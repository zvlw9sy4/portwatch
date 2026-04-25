package notify

// VictorOpsRoutingKey is the routing key that determines which team or
// escalation policy receives the alert in VictorOps / Splunk On-Call.
type VictorOpsConfig struct {
	// WebhookURL is the full REST endpoint URL provided by VictorOps,
	// e.g. https://alert.victorops.com/integrations/generic/…/alert/<api_key>/<routing_key>
	WebhookURL string `yaml:"webhook_url"`

	// MessageType overrides the default message type for opened events.
	// Defaults to "CRITICAL". Closed events always use "RECOVERY".
	MessageType string `yaml:"message_type"`
}

// messageType returns the configured message type, defaulting to CRITICAL.
func (c VictorOpsConfig) messageType() string {
	if c.MessageType == "" {
		return "CRITICAL"
	}
	return c.MessageType
}
