// Package notify provides Notifier implementations for various alerting
// backends. The AmplitudeNotifier sends portwatch events to Amplitude
// Analytics using the HTTP API v2.
//
// # Configuration
//
// Create the notifier with your Amplitude API key:
//
//	n := notify.NewAmplitudeNotifier("YOUR_API_KEY", "")
//
// An empty endpoint string defaults to https://api2.amplitude.com/2/httpapi.
// Override it to point at a proxy or the EU data-residency endpoint:
//
//	n := notify.NewAmplitudeNotifier(key, "https://api.eu.amplitude.com/2/httpapi")
//
// # Events
//
// Each alert.Event is translated to an Amplitude event with:
//   - event_type: "port_opened" or "port_closed"
//   - event_properties.port: the port number
//   - event_properties.protocol: "tcp" or "udp"
//   - user_id: "portwatch" (static sentinel)
package notify
