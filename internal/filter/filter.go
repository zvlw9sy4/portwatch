// Package filter provides port filtering rules for portwatch.
package filter

import "github.com/user/portwatch/internal/scanner"

// Rule defines a single filter rule.
type Rule struct {
	Port     int
	Protocol string // "tcp" or "udp", empty means any
	Comment  string
}

// Filter holds a set of ignore rules.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of rules.
func New(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Apply removes ports that match any ignore rule.
func (f *Filter) Apply(ports []scanner.Port) []scanner.Port {
	if len(f.rules) == 0 {
		return ports
	}
	out := ports[:0:0]
	for _, p := range ports {
		if !f.matches(p) {
			out = append(out, p)
		}
	}
	return out
}

func (f *Filter) matches(p scanner.Port) bool {
	for _, r := range f.rules {
		if r.Port != p.Number {
			continue
		}
		if r.Protocol == "" || r.Protocol == p.Protocol {
			return true
		}
	}
	return false
}
