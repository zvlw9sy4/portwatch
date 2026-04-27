package notify

import (
	"errors"
	"net/url"
)

// HealthChecksConfig holds the configuration fields for the HealthChecks notifier.
// It is intended to be embedded in or referenced by the top-level notify config.
type HealthChecksConfig struct {
	// PingURL is the full URL to POST events to, e.g.
	// https://hc-ping.com/<uuid>/fail
	PingURL string `yaml:"ping_url"`
}

// Validate returns an error if the configuration is not usable.
func (c HealthChecksConfig) Validate() error {
	if c.PingURL == "" {
		return errors.New("healthchecks: ping_url must not be empty")
	}
	u, err := url.ParseRequestURI(c.PingURL)
	if err != nil {
		return errors.New("healthchecks: ping_url is not a valid URL")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("healthchecks: ping_url scheme must be http or https")
	}
	return nil
}

// Build constructs a ready-to-use HealthChecksNotifier from the config.
func (c HealthChecksConfig) Build() (*HealthChecksNotifier, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return NewHealthChecksNotifier(c.PingURL), nil
}
