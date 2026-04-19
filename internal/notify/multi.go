// Package notify provides alert delivery backends.
package notify

import (
	"errors"
	"fmt"

	"github.com/user/portwatch/internal/alert"
)

// MultiNotifier fans out a single Event to multiple Notifier backends,
// collecting all errors rather than stopping at the first failure.
type MultiNotifier struct {
	notifiers []alert.Notifier
}

// NewMultiNotifier creates a MultiNotifier wrapping the given backends.
func NewMultiNotifier(nn ...alert.Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: nn}
}

// Add appends a notifier to the set.
func (m *MultiNotifier) Add(n alert.Notifier) {
	m.notifiers = append(m.notifiers, n)
}

// Notify delivers e to every registered notifier.
// All backends are attempted; a combined error is returned if any fail.
func (m *MultiNotifier) Notify(e alert.Event) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Notify(e); err != nil {
			errs = append(errs, fmt.Errorf("%T: %w", n, err))
		}
	}
	return errors.Join(errs...)
}
