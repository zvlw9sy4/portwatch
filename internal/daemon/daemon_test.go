package daemon

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/state"
)

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestDaemonDetectsOpenPort(t *testing.T) {
	port := freePort(t)

	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("could not bind port: %v", err)
	}
	defer l.Close()

	cfg := config.Default()
	cfg.Ports.Start = port
	cfg.Ports.End = port
	cfg.Interval = 50 * time.Millisecond

	tmpDir := t.TempDir()
	store, _ := state.NewStore(tmpDir + "/state.json")

	received := make(chan struct{}, 1)
	notifier := &captureNotifier{ch: received}
	dispatch := alert.NewDispatcher()
	dispatch.Add(notifier)

	d := New(cfg, store, dispatch)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go d.Run(ctx) //nolint:errcheck

	select {
	case <-received:
		// success
	case <-time.After(400 * time.Millisecond):
		t.Fatal("expected alert not received")
	}
}

type captureNotifier struct {
	ch chan struct{}
}

func (c *captureNotifier) Notify(_ context.Context, events []alert.Event) error {
	if len(events) > 0 {
		select {
		case c.ch <- struct{}{}:
		default:
		}
	}
	return nil
}
