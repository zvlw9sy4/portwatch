// Package notify provides notifier implementations for portwatch alerts.
//
// MatrixNotifier
//
// MatrixNotifier delivers alert events to a Matrix room using the
// Matrix Client-Server API (v3). It requires:
//
//   - homeserver: base URL of the Matrix homeserver (e.g. https://matrix.example.com)
//   - token:      a valid Matrix access token with permission to send messages
//   - roomID:     the fully-qualified Matrix room ID (e.g. !abc123:example.com)
//
// Each alert is sent as an m.room.message event with msgtype m.text.
// A non-2xx HTTP response is treated as an error.
package notify
