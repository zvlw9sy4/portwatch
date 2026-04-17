package daemon_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
	"github.com/user/portwatch/internal/state"
)

func TestDaemonNoAlertsWhenPortsUnchanged(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port

	cfg := config.Default()
	cfg.Ports.Start = port
	cfg.Ports.End = port
	cfg.Interval = 60 * time.Millisecond

	store, _ := state.NewStore(fmt.Sprintf("%s/state.json", t.TempDir()))

	alerts := make(chan struct{}, 10)
	notifier := &countNotifier{ch: alerts}
	dispatch := alert.NewDispatcher()
	dispatch.Add(notifier)

	d := daemon.New(cfg, store, dispatch)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go d.Run(ctx) //nolint:errcheck

	// First tick will alert (new port). Drain it.
	select {
	case <-alerts:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected initial alert")
	}

	// Subsequent ticks should produce no new alerts.
	select {
	case <-alerts:
		t.Fatal("unexpected second alert — port set unchanged")
	case <-ctx.Done():
		// pass
	}
}

type countNotifier struct{ ch chan struct{} }

func (c *countNotifier) Notify(_ context.Context, events []alert.Event) error {
	if len(events) > 0 {
		c.ch <- struct{}{}
	}
	return nil
}
