package relay

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/eskoltech/ethstats/message"
	"github.com/gorilla/websocket"
)

const (
	Api            string = "/api"
	messageHello   string = "hello"
	messagePing    string = "node-ping"
	messageLatency string = "latency"
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
	Secret string
}

// handleRequest is the function to handle all server requests...
func (n *NodeRelay) HandleRequest(w http.ResponseWriter, r *http.Request) {
	nodeConn, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("New Ethereum node connected!")
	go loop(nodeConn, n.Secret)
}

// loop loops as long as the connection is alive and retrieves node packages
func loop(c *websocket.Conn, secret string) {
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
			log.Printf("Received auth message from '%s' node for network %s, node=%s",
				authMsg.ID, authMsg.Info.Network, authMsg.Info.Node)

			// first check if the secret is correct
			if authMsg.Secret != secret {
				log.Printf("Invalid secret from node %s", authMsg.ID)
				return
			}
			sendError := authMsg.SendResponse(c)
			if sendError != nil {
				log.Print(sendError)
			}
			log.Printf("Node '%s' authorized", authMsg.ID)
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
			log.Printf("Server sent pong to node '%s'", ping.ID)
		}

		// Latency from nodes connected to this server.
		if msgType == messageLatency {
			latency, err := msg.GetValue()
			if err != nil {
				log.Print(err)
				return
			}
			log.Printf("Latency: %s", latency)
		}
	}
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
