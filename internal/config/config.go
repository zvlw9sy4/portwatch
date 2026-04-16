package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	Scan     ScanConfig     `yaml:"scan"`
	Alerting AlertingConfig `yaml:"alerting"`
}

type ScanConfig struct {
	Interface string `yaml:"interface"`
	PortStart int    `yaml:"port_start"`
	PortEnd   int    `yaml:"port_end"`
	IntervalS int    `yaml:"interval_seconds"`
}

type AlertingConfig struct {
	LogFile string `yaml:"log_file"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	cfg := Default()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Scan: ScanConfig{
			Interface: "localhost",
			PortStart: 1,
			PortEnd:   65535,
			IntervalS: 60,
		},
		Alerting: AlertingConfig{
			LogFile: "",
		},
	}
}

// Validate checks that the config values are sensible.
func (c *Config) Validate() error {
	if c.Scan.PortStart < 1 || c.Scan.PortStart > 65535 {
		return fmt.Errorf("config: port_start %d out of range", c.Scan.PortStart)
	}
	if c.Scan.PortEnd < 1 || c.Scan.PortEnd > 65535 {
		return fmt.Errorf("config: port_end %d out of range", c.Scan.PortEnd)
	}
	if c.Scan.PortStart > c.Scan.PortEnd {
		return fmt.Errorf("config: port_start %d > port_end %d", c.Scan.PortStart, c.Scan.PortEnd)
	}
	if c.Scan.IntervalS < 1 {
		return fmt.Errorf("config: interval_seconds must be >= 1")
	}
	return nil
}
