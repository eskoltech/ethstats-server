package message

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

// NodeInfo is the information related to the Ethereum node
type NodeInfo struct {
	Name     string `json:"name"`
	Node     string `json:"node"`
	Port     int    `json:"port"`
	Network  string `json:"net"`
	Protocol string `json:"protocol"`
	API      string `json:"api"`
	Os       string `json:"os"`
	OsVer    string `json:"os_v"`
	Client   string `json:"client"`
	History  bool   `json:"canUpdateHistory"`
}

// AuthMessage is the struct sent by the server on the first connection
type AuthMessage struct {
	ID     string   `json:"id"`
	Secret string   `json:"secret"`
	Info   NodeInfo `json:"info"`
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
