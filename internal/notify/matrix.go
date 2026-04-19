package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// MatrixNotifier sends alerts to a Matrix room via the Client-Server API.
type MatrixNotifier struct {
	homeserver string
	token      string
	roomID     string
	client     *http.Client
}

// NewMatrixNotifier returns a MatrixNotifier that posts to the given room.
func NewMatrixNotifier(homeserver, token, roomID string) *MatrixNotifier {
	return &MatrixNotifier{
		homeserver: homeserver,
		token:      token,
		roomID:     roomID,
		client:     &http.Client{},
	}
}

func (m *MatrixNotifier) Notify(e alert.Event) error {
	body := map[string]string{
		"msgtype": "m.text",
		"body":    fmt.Sprintf("[portwatch] %s port %s", e.Kind, e.Port),
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message", m.homeserver, m.roomID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("matrix: unexpected status %d", resp.StatusCode)
	}
	return nil
}
