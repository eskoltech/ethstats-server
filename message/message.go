package message

import (
	"encoding/json"
	"errors"
	"log"
)

// Emit contains the Ethereum message
type Emit struct {
	Content []byte
}

// GetType return the type of the message sent by the Ethereum node
func (e *Emit) GetType() string {
	var content map[string][]interface{}
	err := json.Unmarshal([]byte(e.Content), &content)
	if err != nil {
		log.Println(err)
	}
	result, _ := content["emit"][0].(string)
	return result
}

// Response is the content to send to the Ethereum stats node
func (e *Emit) Response() ([]byte, error) {
	if e.GetType() == "hello" {
		ready := map[string][]interface{}{
			"emit": {"ready"},
		}
		response, err := json.Marshal(ready)
		return response, err
	}
	return nil, errors.New("unknown message type")
}
