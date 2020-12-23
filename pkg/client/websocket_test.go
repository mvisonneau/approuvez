package client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func echoConnectionID(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		if string(message) == "{\"route\":\"get_connection_id\"}" {
			err = c.WriteMessage(mt, []byte("abc123"))
			if err != nil {
				break
			}
		} else {
			err = c.WriteMessage(mt, []byte("error"))
			if err != nil {
				break
			}
		}
	}
}

func TestGetConnectionID(t *testing.T) {
	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echoConnectionID))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	connectionID, err := GetConnectionID(ws)
	assert.Nil(t, err)
	assert.Equal(t, "abc123", connectionID)
}
