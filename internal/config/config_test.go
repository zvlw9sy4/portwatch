package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadValid(t *testing.T) {
	path := writeTempConfig(t, `
scan:
  interface: "127.0.0.1"
  port_start: 1024
  port_end: 9000
  interval_seconds: 30
alerting:
  log_file: "/tmp/portwatch.log"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Scan.PortStart != 1024 {
		t.Errorf("expected port_start 1024, got %d", cfg.Scan.PortStart)
	}
	if cfg.Scan.IntervalS != 30 {
		t.Errorf("expected interval 30, got %d", cfg.Scan.IntervalS)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidatePortRange(t *testing.T) {
	cfg := config.Default()
	cfg.Scan.PortStart = 9000
	cfg.Scan.PortEnd = 1000
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for inverted port range")
	}
}

func TestValidateInterval(t *testing.T) {
	cfg := config.Default()
	cfg.Scan.IntervalS = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for zero interval")
	}
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
}
