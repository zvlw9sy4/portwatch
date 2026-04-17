// Package metrics tracks runtime statistics for the portwatch daemon.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of daemon statistics.
type Snapshot struct {
	ScansTotal   int64
	AlertsTotal  int64
	LastScanAt   time.Time
	LastAlertAt  time.Time
	UptimeSince  time.Time
}

// Collector accumulates runtime metrics in a thread-safe manner.
type Collector struct {
	mu          sync.RWMutex
	scansTotal  int64
	alertsTotal int64
	lastScanAt  time.Time
	lastAlertAt time.Time
	start       time.Time
}

// NewCollector returns a Collector initialised with the current time.
func NewCollector() *Collector {
	return &Collector{start: time.Now()}
}

// RecordScan increments the scan counter and records the timestamp.
func (c *Collector) RecordScan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.lastScanAt = time.Now()
}

// RecordAlert increments the alert counter and records the timestamp.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal++
	c.lastAlertAt = time.Now()
}

// Snapshot returns a copy of the current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		ScansTotal:  c.scansTotal,
		AlertsTotal: c.alertsTotal,
		LastScanAt:  c.lastScanAt,
		LastAlertAt: c.lastAlertAt,
		UptimeSince: c.start,
	}
}
