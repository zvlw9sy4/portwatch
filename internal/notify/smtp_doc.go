// Package notify provides notifier implementations for portwatch alert events.
//
// SMTPNotifier
//
// SMTPNotifier delivers alert events as plain-text email messages using the
// standard net/smtp package. Both STARTTLS (the default, port 587) and
// implicit TLS (port 465) modes are supported via the UseTLS option.
//
// Basic usage:
//
//	n := notify.NewSMTPNotifier(notify.SMTPOptions{
//		Host:     "smtp.example.com",
//		Port:     587,
//		Username: "alerts@example.com",
//		Password: "secret",
//		From:     "alerts@example.com",
//		To:       []string{"ops@example.com"},
//	})
//
// For implicit TLS on port 465 set UseTLS: true.
//
// If the events slice is empty, Notify is a no-op and no connection is made.
package notify
