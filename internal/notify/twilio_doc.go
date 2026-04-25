// Package notify provides notifier implementations for portwatch alerts.
//
// # Twilio SMS Notifier
//
// TwilioNotifier delivers port-change alerts as SMS messages using the
// Twilio Programmable Messaging REST API.
//
// # Configuration
//
// The notifier requires four parameters:
//
//	- accountSID  — Twilio Account SID (starts with "AC")
//	- authToken   — Twilio Auth Token
//	- from        — Sender phone number in E.164 format (+15551234567)
//	- to          — Recipient phone number in E.164 format
//
// # Usage
//
//	n := notify.NewTwilioNotifier(sid, token, "+15550001111", "+15559998888")
//	multi.Add(n)
//
// Each scan cycle that produces events will result in a single SMS
// containing all opened/closed port changes, one per line.
package notify
