package alert

import (
	"log"
)

// Dispatcher holds a list of Notifiers and fans out Events to all of them.
type Dispatcher struct {
	notifiers []Notifier
}

// NewDispatcher creates a Dispatcher with the provided notifiers.
func NewDispatcher(notifiers ...Notifier) *Dispatcher {
	return &Dispatcher{notifiers: notifiers}
}

// Add registers an additional Notifier.
func (d *Dispatcher) Add(n Notifier) {
	d.notifiers = append(d.notifiers, n)
}

// Dispatch sends all events to every registered Notifier.
// Errors are logged but do not stop delivery to remaining notifiers.
func (d *Dispatcher) Dispatch(events []Event) {
	for _, e := range events {
		for _, n := range d.notifiers {
			if err := n.Notify(e); err != nil {
				log.Printf("portwatch: notifier error: %v", err)
			}
		}
	}
}
