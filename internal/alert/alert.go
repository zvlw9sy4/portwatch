package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a port change alert.
type Event struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      scanner.Port
}

// Notifier sends alert events to a destination.
type Notifier interface {
	Notify(e Event) error
}

// LogNotifier writes alerts to an io.Writer.
type LogNotifier struct {
	Out io.Writer
}

// NewLogNotifier returns a LogNotifier writing to stdout by default.
func NewLogNotifier(out io.Writer) *LogNotifier {
	if out == nil {
		out = os.Stdout
	}
	return &LogNotifier{Out: out}
}

// Notify formats and writes the event.
func (l *LogNotifier) Notify(e Event) error {
	_, err := fmt.Fprintf(
		l.Out,
		"[%s] %s %s — %s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Port.String(),
		e.Message,
	)
	return err
}

// NotifyAll sends all events to the notifier, returning the first error
// encountered along with the number of events successfully delivered.
func NotifyAll(n Notifier, events []Event) (int, error) {
	for i, e := range events {
		if err := n.Notify(e); err != nil {
			return i, fmt.Errorf("alert: failed to notify event %d: %w", i, err)
		}
	}
	return len(events), nil
}

// BuildEvents converts a diff result into a slice of alert Events.
func BuildEvents(opened, closed []scanner.Port) []Event {
	now := time.Now()
	events := make([]Event, 0, len(opened)+len(closed))
	for _, p := range opened {
		events = append(events, Event{
			Timestamp: now,
			Level:     LevelAlert,
			Message:   "port newly opened",
			Port:      p,
		})
	}
	for _, p := range closed {
		events = append(events, Event{
			Timestamp: now,
			Level:     LevelWarn,
			Message:   "port closed",
			Port:      p,
		})
	}
	return events
}
