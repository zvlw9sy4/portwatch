package notify

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// SMTPNotifier sends alert events via SMTP with optional TLS support.
type SMTPNotifier struct {
	host     string
	port     int
	username string
	password string
	from     string
	to       []string
	useTLS   bool
}

// SMTPOptions configures the SMTP notifier.
type SMTPOptions struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
	UseTLS   bool
}

// NewSMTPNotifier creates a new SMTPNotifier with the given options.
func NewSMTPNotifier(opts SMTPOptions) *SMTPNotifier {
	if opts.Port == 0 {
		opts.Port = 587
	}
	return &SMTPNotifier{
		host:     opts.Host,
		port:     opts.Port,
		username: opts.Username,
		password: opts.Password,
		from:     opts.From,
		to:       opts.To,
		useTLS:   opts.UseTLS,
	}
}

// Notify sends an email for each alert event.
func (s *SMTPNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, e := range events {
		sb.WriteString(fmt.Sprintf("[%s] port %s/%s\n", e.Kind, e.Port.Number, e.Port.Proto))
	}

	subject := fmt.Sprintf("portwatch: %d port change(s) detected", len(events))
	body := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s",
		strings.Join(s.to, ", "), s.from, subject, sb.String())

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	if s.useTLS {
		return s.sendTLS(addr, auth, []byte(body))
	}
	return smtp.SendMail(addr, auth, s.from, s.to, []byte(body))
}

func (s *SMTPNotifier) sendTLS(addr string, auth smtp.Auth, body []byte) error {
	tlsCfg := &tls.Config{ServerName: s.host}
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return fmt.Errorf("smtp tls dial: %w", err)
	}
	host, _, _ := net.SplitHostPort(addr)
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp new client: %w", err)
	}
	defer c.Quit() //nolint:errcheck
	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := c.Mail(s.from); err != nil {
		return err
	}
	for _, r := range s.to {
		if err := c.Rcpt(r); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = w.Write(body)
	return err
}
