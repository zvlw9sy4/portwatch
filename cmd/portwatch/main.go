package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
	"github.com/user/portwatch/internal/state"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	statePath := flag.String("state", "/tmp/portwatch_state.json", "path to state file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Printf("[main] config load failed, using defaults: %v", err)
		cfg = config.Default()
	}

	store, err := state.NewStore(*statePath)
	if err != nil {
		log.Fatalf("[main] state store: %v", err)
	}

	dispatch := alert.NewDispatcher()
	dispatch.Add(alert.NewLogNotifier(os.Stdout))

	d := daemon.New(cfg, store, dispatch)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("[main] portwatch started (ports %d-%d, interval %s)",
		cfg.Ports.Start, cfg.Ports.End, cfg.Interval)

	if err := d.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("[main] daemon error: %v", err)
	}
}
