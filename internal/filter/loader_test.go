package filter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/filter"
)

func writeTempFilter(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "filter.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadRulesValid(t *testing.T) {
	p := writeTempFilter(t, "ignore:\n  - port: 22\n    protocol: tcp\n    comment: ssh\n")
	rules, err := filter.LoadRules(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 || rules[0].Port != 22 {
		t.Fatalf("unexpected rules: %v", rules)
	}
}

func TestLoadRulesMissingFile(t *testing.T) {
	_, err := filter.LoadRules("/nonexistent/filter.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadRulesInvalidPort(t *testing.T) {
	p := writeTempFilter(t, "ignore:\n  - port: 99999\n")
	_, err := filter.LoadRules(p)
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestLoadRulesEmpty(t *testing.T) {
	p := writeTempFilter(t, "ignore: []\n")
	rules, err := filter.LoadRules(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 0 {
		t.Fatalf("expected empty rules")
	}
}

func TestLoadRulesInvalidYAML(t *testing.T) {
	p := writeTempFilter(t, "ignore:\n  - port: [unclosed\n")
	_, err := filter.LoadRules(p)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
