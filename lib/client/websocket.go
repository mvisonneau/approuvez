package client

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// GetConnectionID returns the ID of the currently established websocket connection with API Gateway
func (c *Client) GetConnectionID() (string, error) {
	log.Debugf("Fetching connection ID from websocket")
	if err := c.Websocket.WriteMessage(websocket.TextMessage, []byte("{\"route\":\"get_connection_id\"}")); err != nil {
		return "", err
	}

	log.Debugf("Request for connection ID sent, waiting for response on websocket")
	_, connectionID, err := c.Websocket.ReadMessage()
	if err != nil {
		return "", err
	}

	log.Infof("Connection ID: %s", string(connectionID))
	return string(connectionID), nil
}
