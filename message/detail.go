package message

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

// AuthMessage is the struct sent by the server on the first connection
type AuthMessage struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// SendResponse send the ready response to the node to initiate the communication
func (a *AuthMessage) SendResponse(c *websocket.Conn) error {
	ready := map[string][]interface{}{"emit": {"ready"}}
	response, err := json.Marshal(ready)
	if err != nil {
		return err
	}
	err = c.WriteMessage(1, response)
	if err != nil {
		return err
	}
	return nil
}
