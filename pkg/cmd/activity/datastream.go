package activity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/spancli/pkg/global"
)

type websocketAuth struct {
	JWT          string `json:"jwt"`
	CollectionID string `json:"collectionId"`
	DeviceID     string `json:"deviceId"`
	GatewayID    string `json:"gatewayId"`
}

// NewActivityEventStream creates a live activity event stream from the API
func NewActivityEventStream(token string, jwt string, collectionID string, deviceID string, gatewayID string) (*ActivityEventStream, error) {
	wsURL := fmt.Sprintf("wss://api.lab5e.com/span/collections/%s/activity", collectionID)
	if global.Options.OverrideEndpoint != "" {
		u, err := url.Parse(global.Options.OverrideEndpoint)
		if err != nil {
			return nil, err
		}

		wsURL = fmt.Sprintf("ws://%s/span/collections/%s/activity", u.Host, collectionID)
	}

	header := http.Header{}
	if token != "" {
		header.Add("X-API-Token", token)
	}

	dialer := websocket.Dialer{}
	ws, _, err := dialer.Dial(wsURL, header)
	if err != nil {
		return nil, fmt.Errorf("error dialing websocket: %v", err)
	}

	if token == "" && jwt != "" {
		if err := ws.WriteJSON(websocketAuth{
			JWT:          jwt,
			CollectionID: collectionID,
			DeviceID:     deviceID,
			GatewayID:    gatewayID,
		}); err != nil {
			ws.Close()
			return nil, err
		}
	}
	return &ActivityEventStream{ws}, nil
}

type ActivityEventStream struct {
	ws *websocket.Conn
}

func (d *ActivityEventStream) Recv() (spanapi.ActivityEvent, error) {
	_, msgBytes, err := d.ws.ReadMessage()
	if err != nil {
		return spanapi.ActivityEvent{}, err
	}

	m := spanapi.ActivityEvent{}
	err = json.Unmarshal(msgBytes, &m)
	if err != nil {
		return spanapi.ActivityEvent{}, err
	}

	return m, nil
}

func (d *ActivityEventStream) Close() error {
	return d.ws.Close()
}
