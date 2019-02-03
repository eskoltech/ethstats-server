package message

import (
	"encoding/json"
	"log"
)

// Message contains the Ethereum message
type Message struct {
	Content []byte
}

// GetType return the type of the message sent by the Ethereum node
func (e *Message) GetType() string {
	var content map[string][]interface{}
	err := json.Unmarshal([]byte(e.Content), &content)
	if err != nil {
		log.Println(err)
	}
	result, _ := content["emit"][0].(string)
	return result
}

func (e *Message) GetValue() []byte {
	var content map[string][]interface{}
	err := json.Unmarshal([]byte(e.Content), &content)
	if err != nil {
		log.Println(err)
	}
	result, _ := content["emit"][1].(interface{})
	val, _ := json.Marshal(result)
	return val
}
