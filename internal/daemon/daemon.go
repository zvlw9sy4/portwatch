package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Daemon orchestrates periodic port scanning and alerting.
type Daemon struct {
	cfg      *config.Config
	scanner  *scanner.Scanner
	store    *state.Store
	dispatch *alert.Dispatcher
}

// New creates a Daemon from the provided config.
func New(cfg *config.Config, store *state.Store, dispatch *alert.Dispatcher) *Daemon {
	s := scanner.NewScanner(cfg.Ports.Start, cfg.Ports.End, cfg.Timeout)
	return &Daemon{cfg: cfg, scanner: s, store: store, dispatch: dispatch}
}

// Run starts the scan loop, blocking until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	// Run an immediate first scan.
	if err := d.tick(ctx); err != nil {
		log.Printf("[daemon] initial scan error: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := d.tick(ctx); err != nil {
				log.Printf("[daemon] scan error: %v", err)
			}
		case <-ctx.Done():
			log.Println("[daemon] shutting down")
			return ctx.Err()
		}
	}
}

func (d *Daemon) tick(ctx context.Context) error {
	current, err := d.scanner.Scan(ctx)
	if err != nil {
		return err
	}

	previous, _ := d.store.Load()
	diff := scanner.Compare(previous, current)

	if len(diff.Opened)+len(diff.Closed) > 0 {
		events := alert.BuildEvents(diff)
		d.dispatch.Send(ctx, events)
	}

	return d.store.Save(current)
}
