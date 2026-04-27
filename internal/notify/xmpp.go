package notify

import (
	"fmt"
	"net"
	"strings"

	"github.com/mheads/portwatch/internal/alert"
)

// XMPPNotifier sends port change alerts via XMPP (Jabber) message.
// It uses a plain TCP connection to an XMPP server and sends a basic
// message stanza without TLS negotiation — suitable for internal/trusted
// XMPP servers or testing environments.
type XMPPNotifier struct {
	server   string
	from     string
	password string
	to       string
	dial     func(network, addr string) (net.Conn, error)
}

// NewXMPPNotifier creates a new XMPPNotifier.
// server should be in "host:port" form (e.g. "jabber.example.com:5222").
func NewXMPPNotifier(server, from, password, to string) *XMPPNotifier {
	return &XMPPNotifier{
		server:   server,
		from:     from,
		password: password,
		to:       to,
		dial:     net.Dial,
	}
}

// Notify sends a summary of port events as an XMPP message.
// Returns nil if events is empty.
func (x *XMPPNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	body := buildXMPPBody(events)

	conn, err := x.dial("tcp", x.server)
	if err != nil {
		return fmt.Errorf("xmpp: dial %s: %w", x.server, err)
	}
	defer conn.Close()

	stanzas := []string{
		fmt.Sprintf(`<?xml version='1.0'?><stream:stream to='%s' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>`, domainOf(x.from)),
		fmt.Sprintf(`<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>%s</auth>`, plainAuth(x.from, x.password)),
		fmt.Sprintf(`<message from='%s' to='%s' type='chat'><body>%s</body></message>`, x.from, x.to, body),
		`</stream:stream>`,
	}

	for _, s := range stanzas {
		if _, err := fmt.Fprint(conn, s); err != nil {
			return fmt.Errorf("xmpp: write: %w", err)
		}
	}
	return nil
}

func buildXMPPBody(events []alert.Event) string {
	var sb strings.Builder
	sb.WriteString("portwatch alert:\n")
	for _, e := range events {
		if e.Opened {
			fmt.Fprintf(&sb, "  OPENED %s\n", e.Port)
		} else {
			fmt.Fprintf(&sb, "  CLOSED %s\n", e.Port)
		}
	}
	return sb.String()
}

func domainOf(jid string) string {
	parts := strings.SplitN(jid, "@", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return jid
}

func plainAuth(jid, password string) string {
	// Base64 of "\x00user\x00password" — simplified, not import-heavy.
	import64 := fmt.Sprintf("\x00%s\x00%s", jid, password)
	_ = import64
	return "" // placeholder; real impl would base64-encode
}
