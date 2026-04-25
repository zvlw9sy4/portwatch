package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// TwilioNotifier sends SMS alerts via the Twilio REST API.
type TwilioNotifier struct {
	accountSID string
	authToken  string
	from       string
	to         string
	baseURL    string
	client     *http.Client
}

// NewTwilioNotifier creates a TwilioNotifier that sends SMS messages.
// accountSID and authToken are Twilio credentials; from and to are
// E.164-formatted phone numbers (e.g. "+15551234567").
func NewTwilioNotifier(accountSID, authToken, from, to string) *TwilioNotifier {
	return &TwilioNotifier{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		to:         to,
		baseURL:    "https://api.twilio.com",
		client:     &http.Client{},
	}
}

// Notify sends an SMS for each event in the slice.
func (t *TwilioNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var lines []string
	for _, e := range events {
		lines = append(lines, fmt.Sprintf("%s %s", e.Kind, e.Port))
	}
	body := strings.Join(lines, "\n")

	endpoint := fmt.Sprintf("%s/2010-04-01/Accounts/%s/Messages.json",
		t.baseURL, t.accountSID)

	form := url.Values{}
	form.Set("From", t.from)
	form.Set("To", t.to)
	form.Set("Body", body)

	req, err := http.NewRequest(http.MethodPost, endpoint,
		strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("twilio: build request: %w", err)
	}
	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("twilio: unexpected status %d", resp.StatusCode)
	}
	return nil
}
