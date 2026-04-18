package notify

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// EmailConfig holds SMTP connection and addressing details.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// EmailNotifier sends alert events via SMTP.
type EmailNotifier struct {
	cfg  EmailConfig
	auth smtp.Auth
}

// NewEmailNotifier constructs an EmailNotifier and configures PLAIN auth.
func NewEmailNotifier(cfg EmailConfig) *EmailNotifier {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &EmailNotifier{cfg: cfg, auth: auth}
}

// Notify sends a single alert event as an email.
func (e *EmailNotifier) Notify(ev alert.Event) error {
	subject := fmt.Sprintf("[portwatch] port %s %s", ev.Port, ev.Kind)
	body := fmt.Sprintf("Port: %s\nProtocol: %s\nEvent: %s\n",
		ev.Port, ev.Protocol, ev.Kind)

	msg := []byte("To: " + strings.Join(e.cfg.To, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	return smtp.SendMail(addr, e.auth, e.cfg.From, e.cfg.To, msg)
}
