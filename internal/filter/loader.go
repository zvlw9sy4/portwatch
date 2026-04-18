package filter

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RuleConfig is the YAML representation of a filter rule.
type RuleConfig struct {
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	Comment  string `yaml:"comment"`
}

// LoadRules reads filter rules from a YAML file.
// The file should contain a top-level "ignore" list.
func LoadRules(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("filter: read %s: %w", path, err)
	}
	var cfg struct {
		Ignore []RuleConfig `yaml:"ignore"`
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("filter: parse %s: %w", path, err)
	}
	rules := make([]Rule, 0, len(cfg.Ignore))
	for _, rc := range cfg.Ignore {
		if rc.Port <= 0 || rc.Port > 65535 {
			return nil, fmt.Errorf("filter: invalid port %d", rc.Port)
		}
		rules = append(rules, Rule{Port: rc.Port, Protocol: rc.Protocol, Comment: rc.Comment})
	}
	return rules, nil
}
