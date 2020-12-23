package client

import (
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// GetConnectionID returns the ID of the currently established websocket connection with API Gateway
func GetConnectionID(c *websocket.Conn) (string, error) {
	log.Debug("fetching connection ID from websocket")
	if err := c.WriteMessage(websocket.TextMessage, []byte("{\"route\":\"get_connection_id\"}")); err != nil {
		return "", err
	}

	log.Debug("request for connection ID sent, waiting for response on websocket")
	_, connectionID, err := c.ReadMessage()
	if err != nil {
		return "", err
	}

	log.WithFields(
		log.Fields{
			"connection-id": string(connectionID),
		},
	).Info("retrieved connection ID")
	return string(connectionID), nil
}

// KeepAlive make simple pings upon a WebSocket connection to attempt keeping it alive
func KeepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout / 2)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				log.Warn("websocket keepalive routine stopped")
				return
			case <-ticker.C:
				if err := c.WriteMessage(websocket.TextMessage, []byte("keepalive")); err != nil {
					log.Error(err.Error())
					return
				}
			}
		}
	}()
}
