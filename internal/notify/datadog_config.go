package notify

// DatadogConfig holds configuration for the Datadog notifier.
// It is intended to be populated from the portwatch YAML config file.
type DatadogConfig struct {
	// APIKey is the Datadog API key (required).
	APIKey string `yaml:"api_key"`

	// BaseURL overrides the default Datadog API endpoint.
	// Leave empty to use https://api.datadoghq.com (US region).
	// EU customers should set https://api.datadoghq.eu.
	BaseURL string `yaml:"base_url"`
}

// IsEnabled returns true when the DatadogConfig has an API key configured.
func (c DatadogConfig) IsEnabled() bool {
	return c.APIKey != ""
}

// Build constructs a DatadogNotifier from the configuration.
func (c DatadogConfig) Build() *DatadogNotifier {
	return NewDatadogNotifier(c.APIKey, c.BaseURL)
}
