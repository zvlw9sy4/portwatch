package notify

// DatadogConfig holds configuration for the Datadog event notifier.
type DatadogConfig struct {
	// APIKey is the Datadog API key.
	APIKey string `yaml:"api_key"`

	// Site is the Datadog intake site, e.g. "datadoghq.com" (default) or
	// "datadoghq.eu" for the EU region.
	Site string `yaml:"site"`

	// Tags is an optional list of tags appended to every event.
	Tags []string `yaml:"tags"`

	// AlertType overrides the default alert type for opened-port events.
	// Valid values: "error", "warning", "info", "success". Defaults to "warning".
	AlertType string `yaml:"alert_type"`
}

// baseURL returns the Datadog API base URL for the configured site.
func (c DatadogConfig) baseURL() string {
	site := c.Site
	if site == "" {
		site = "datadoghq.com"
	}
	return "https://api." + site
}

// alertType returns the configured alert type, defaulting to "warning".
func (c DatadogConfig) alertType() string {
	if c.AlertType == "" {
		return "warning"
	}
	return c.AlertType
}
