package relay

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/eskoltech/ethstats/message"
	"github.com/eskoltech/ethstats/service"
	"github.com/gorilla/websocket"
)

const (
	Api            string = "/api"
	messageHello   string = "hello"
	messagePing    string = "node-ping"
	messageLatency string = "latency"
	messageBlock   string = "block"
	messageHistory string = "history"
	messagePending string = "pending"
	messageStats   string = "stats"
)

// upgradeConnection upgrade only HTTP request to the /api endpoint
var upgradeConnection = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return strings.Compare(r.RequestURI, Api) == 0
	},
}

// NodeRelay contains the secret used to authenticate the communication between
// the Ethereum node and this server
type NodeRelay struct {
	Secret  string
	Service *service.Channel
}

// HandleRequest is the function to handle all server requests that came from
// Ethereum nodes
func (n *NodeRelay) HandleRequest(w http.ResponseWriter, r *http.Request) {
	nodeConn, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("New Ethereum node connected!")
	go n.loop(nodeConn)
}

// loop loops as long as the connection is alive and retrieves node packages
func (n *NodeRelay) loop(c *websocket.Conn) {
	// Close connection if an unexpected error occurs
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(c)
	// Client loop
	for {
		_, content, err := c.ReadMessage()
		if err != nil {
			break
		}
		// Create emitted message from the node
		msg := message.Message{Content: content}
		msgType, err := msg.GetType()
		if err != nil {
			log.Print(err)
			return
		}

		// If message type is hello, we need to check if the secret is
		// correct, and then, send a ready message
		if msgType == messageHello {
			// Get value from JSON to store it and process it later to calculate
			// node latency etc
			authMsg, parseError := parseAuthMessage(msg)
			if parseError != nil {
				log.Print(parseError)
				return
			}
			// first check if the secret is correct
			if authMsg.Secret != n.Secret {
				log.Printf("Invalid secret from node %s", authMsg.ID)
				return
			}
			sendError := authMsg.SendResponse(c)
			if sendError != nil {
				log.Print(sendError)
			}
			go func(s *service.Channel, a []byte) {
				s.Message <- a
			}(n.Service, content)
		}

		// When the node emit a ping message, we need to respond with pong
		// before five seconds to authorize that node to sent reports
		if msgType == messagePing {
			ping, err := parseNodePingMessage(msg)
			if err != nil {
				log.Print(err)
				return
			}
			sendError := ping.SendResponse(c)
			if sendError != nil {
				log.Print(sendError)
			}
			go func(s *service.Channel, p []byte) {
				s.Message <- p
			}(n.Service, content)
		}

		// Send the content sent by the nodes directly to the consumer clients.
		// Only message types recognized by this server
		if isValidMessage(msgType) {
			go func(s *service.Channel, l []byte) {
				s.Message <- l
			}(n.Service, content)
		}
	}
}

// isValidMessage return true if the message type is know, otherwise return false
func isValidMessage(msgType string) bool {
	return msgType == messageLatency || msgType == messageBlock || msgType == messageHistory || msgType == messagePending || msgType == messageStats
}

// parseNodePingMessage parse the current ping message sent bu the Ethereum node
// and creates a message.NodePing struct with that info
func parseNodePingMessage(msg message.Message) (*message.NodePing, error) {
	value, err := msg.GetValue()
	if err != nil {
		return &message.NodePing{}, err
	}
	var ping message.NodePing
	err = json.Unmarshal(value, &ping)
	return &ping, err
}

// parseAuthMessage parse the current byte array and transforms it to an AuthMessage struct.
// If an error occurs when json unmarshal, an error is returned
func parseAuthMessage(msg message.Message) (*message.AuthMessage, error) {
	value, err := msg.GetValue()
	if err != nil {
		return &message.AuthMessage{}, err
	}
	var detail message.AuthMessage
	err = json.Unmarshal(value, &detail)
	return &detail, err
}
