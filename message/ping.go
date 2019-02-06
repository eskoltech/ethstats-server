package message

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

// NodePing contains the last time the node is alive
type NodePing struct {
	ID   string `json:"id"`
	Time string `json:"clientTime"`
}

// SendResponse send the pong response to the node
func (n *NodePing) SendResponse(c *websocket.Conn) error {
	ready := map[string][]interface{}{"emit": {"node-pong", n.ID}}
	response, err := json.Marshal(ready)
	if err != nil {
		return err
	}
	// message type is always 1
	err = c.WriteMessage(1, response)
	if err != nil {
		return err
	}
	return nil
}
