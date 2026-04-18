package watch_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

const validConfig = `
port_range: "1-1024"
interval: 5
state_file: "/tmp/portwatch.json"
`

const updatedConfig = `
port_range: "1-2048"
interval: 10
state_file: "/tmp/portwatch.json"
`

func TestReloadableConfigGet(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/config.yaml"
	_ = os.WriteFile(path, []byte(validConfig), 0o644)

	rc, w, err := watch.NewReloadableConfig(path, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Start()
	defer w.Stop()

	cfg := rc.Get()
	if cfg.PortRange != "1-1024" {
		t.Errorf("expected port_range 1-1024, got %s", cfg.PortRange)
	}
}

func TestReloadableConfigReloadsOnChange(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/config.yaml"
	_ = os.WriteFile(path, []byte(validConfig), 0o644)

	rc, w, err := watch.NewReloadableConfig(path, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	_ = os.WriteFile(path, []byte(updatedConfig), 0o644)
	time.Sleep(80 * time.Millisecond)

	cfg := rc.Get()
	if cfg.PortRange != "1-2048" {
		t.Errorf("expected reloaded port_range 1-2048, got %s", cfg.PortRange)
	}
}

func TestNewReloadableConfigMissingFile(t *testing.T) {
	_, _, err := watch.NewReloadableConfig("/no/such/file.yaml", 1)
	if err == nil {
		t.Error("expected error for missing config file")
	}
}
