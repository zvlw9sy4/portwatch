package notify

// TwilioConfig holds the configuration fields for the Twilio SMS notifier.
// It is intended to be embedded in or deserialized from the portwatch
// YAML configuration file under a "twilio" key.
//
//	twilio:
//	  account_sid: ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//	  auth_token:  your_auth_token
//	  from:        "+15550001111"
//	  to:          "+15559998888"
type TwilioConfig struct {
	// AccountSID is the Twilio Account SID.
	AccountSID string `yaml:"account_sid"`
	// AuthToken is the Twilio Auth Token.
	AuthToken string `yaml:"auth_token"`
	// From is the sender phone number in E.164 format.
	From string `yaml:"from"`
	// To is the recipient phone number in E.164 format.
	To string `yaml:"to"`
}

// IsConfigured reports whether all required Twilio fields are present.
func (c TwilioConfig) IsConfigured() bool {
	return c.AccountSID != "" && c.AuthToken != "" && c.From != "" && c.To != ""
}

// Build constructs a TwilioNotifier from the config. It panics if the
// config is incomplete; callers should check IsConfigured first.
func (c TwilioConfig) Build() *TwilioNotifier {
	if !c.IsConfigured() {
		panic("notify: TwilioConfig.Build called with incomplete config")
	}
	return NewTwilioNotifier(c.AccountSID, c.AuthToken, c.From, c.To)
}
