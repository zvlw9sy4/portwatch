package notify

import (
	"fmt"
	"log/syslog"

	"github.com/user/portwatch/internal/alert"
)

// SyslogNotifier sends alerts to the local syslog daemon.
type SyslogNotifier struct {
	writer *syslog.Writer
	tag    string
}

// NewSyslogNotifier creates a SyslogNotifier using the given tag.
// priority defaults to LOG_ALERT | LOG_DAEMON.
func NewSyslogNotifier(tag string) (*SyslogNotifier, error) {
	w, err := syslog.New(syslog.LOG_ALERT|syslog.LOG_DAEMON, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: open: %w", err)
	}
	return &SyslogNotifier{writer: w, tag: tag}, nil
}

// Notify sends a single event to syslog.
func (s *SyslogNotifier) Notify(e alert.Event) error {
	msg := fmt.Sprintf("portwatch [%s] port %d/%s — %s",
		e.Kind, e.Port.Number, e.Port.Protocol, e.Port.Address)
	switch e.Kind {
	case alert.EventOpened:
		return s.writer.Alert(msg)
	default:
		return s.writer.Info(msg)
	}
}

// Close releases the underlying syslog connection.
func (s *SyslogNotifier) Close() error {
	return s.writer.Close()
}
