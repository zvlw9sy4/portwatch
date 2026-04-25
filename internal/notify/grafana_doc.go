// Package notify provides notifier implementations for portwatch.
//
// # Grafana Notifier
//
// GrafanaNotifier posts an annotation to a Grafana instance for every
// port-change event detected by portwatch. This allows operators to
// correlate network-topology changes with dashboard time-series data.
//
// # Configuration
//
// Construct a GrafanaNotifier with the annotations API endpoint, a
// service-account API key, and an optional list of tags:
//
//	n := notify.NewGrafanaNotifier(
//		"http://grafana.example.com/api/annotations",
//		"glsa_xxxxxxxxxxxxxxxxxxxx",
//		[]string{"portwatch", "production"},
//	)
//
// When tags is nil the notifier defaults to ["portwatch"].
//
// Each event generates a separate annotation request so that individual
// port transitions appear as distinct markers on the Grafana timeline.
package notify
