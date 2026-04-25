package notify

// OpsGenieConfig holds configuration for the OpsGenie notifier.
type OpsGenieConfig struct {
	// APIKey is the OpsGenie API key used for authentication.
	APIKey string `yaml:"api_key"`

	// Region controls which OpsGenie API endpoint is used.
	// Valid values are "us" (default) and "eu".
	Region string `yaml:"region"`

	// Priority is the alert priority (P1–P5). Defaults to P3.
	Priority string `yaml:"priority"`

	// Tags is an optional list of tags to attach to every alert.
	Tags []string `yaml:"tags"`
}

// baseURL returns the correct OpsGenie API base URL for the configured region.
func (c OpsGenieConfig) baseURL() string {
	if c.Region == "eu" {
		return "https://api.eu.opsgenie.com"
	}
	return "https://api.opsgenie.com"
}

// priority returns the configured priority, defaulting to P3.
func (c OpsGenieConfig) priority() string {
	if c.Priority == "" {
		return "P3"
	}
	return c.Priority
}
