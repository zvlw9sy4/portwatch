package notify

import (
	"fmt"
)

// VictorOpsConfig holds configuration for the VictorOps notifier.
type VictorOpsConfig struct {
	// EndpointURL is the base REST endpoint URL provided by VictorOps.
	EndpointURL string `yaml:"endpoint_url"`
	// RoutingKey identifies the escalation policy to trigger.
	RoutingKey string `yaml:"routing_key"`
}

// Validate returns an error if the configuration is incomplete.
func (c VictorOpsConfig) Validate() error {
	if c.EndpointURL == "" {
		return fmt.Errorf("victorops: endpoint_url is required")
	}
	if c.RoutingKey == "" {
		return fmt.Errorf("victorops: routing_key is required")
	}
	return nil
}

// Build constructs a VictorOpsNotifier from the configuration.
func (c VictorOpsConfig) Build() (*VictorOpsNotifier, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return NewVictorOpsNotifier(c.EndpointURL, c.RoutingKey), nil
}
