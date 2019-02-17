package relay

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/eskoltech/ethstats-server/message"
	"github.com/eskoltech/ethstats-server/service"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	// Api is the public endpoint used to send stats from nodes to this server
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
	secret  string
	service *service.Channel
}

// New creates a new NodeRelay struct with required fields
func New(service *service.Channel, secret string) *NodeRelay {
	defer func() { log.Info("Node relay started successfully") }()
	return &NodeRelay{
		service: service,
		secret:  secret,
	}
}

// Close closes the connection between this server and all Ethereum nodes connected to it
func (n *NodeRelay) Close() {
	log.Info("Prepared to close connection with nodes...")
	close(n.service.Message)
}

// HandleRequest is the function to handle all server requests that came from
// Ethereum nodes
func (n *NodeRelay) HandleRequest(w http.ResponseWriter, r *http.Request) {
	nodeConn, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Warningf("Error establishing node connection: %s", err)
		return
	}
	log.Infof("New Ethereum node connected! (addr=%s, host=%s)", r.RemoteAddr, r.Host)
	go n.loop(nodeConn)
}

// loop loops as long as the connection is alive and retrieves node packages
func (n *NodeRelay) loop(c *websocket.Conn) {
	// Close connection if an unexpected error occurs and delete the node
	// from the map of connected nodes...
	defer func(conn *websocket.Conn) {
		delete(n.service.Nodes, conn.RemoteAddr().String())
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Warningf("Connection with node closed, there are %d connected nodes", len(n.service.Nodes))
	}(c)
	// Client loop
	for {
		_, content, err := c.ReadMessage()
		if err != nil {
			log.Errorf("Error reading message from client, %s", err)
			break
		}
		// Create emitted message from the node
		msg := message.Message{Content: content}
		msgType, err := msg.GetType()
		if err != nil {
			log.Warningf("Can't get type of message sent by the node: %s", err)
			return
		}

		// If message type is hello, we need to check if the secret is
		// correct, and then, send a ready message
		if msgType == messageHello {
			// Get value from JSON to store it and process it later to calculate
			// node latency etc
			authMsg, parseError := parseAuthMessage(msg)
			if parseError != nil {
				log.Warningf("Can't parse authorization message sent by node[%s], error: %s", authMsg.ID, parseError)
				return
			}
			// first check if the secret is correct
			if authMsg.Secret != n.secret {
				log.Errorf("Invalid secret from node %s, can't get stats", authMsg.ID)
				return
			}
			sendError := authMsg.SendResponse(c)
			if sendError != nil {
				log.Errorf("Error sending authorization response to node[%s], error: %s", authMsg.ID, sendError)
				return
			}
			n.service.Message <- content

			// use node addr as identifier to check node availability
			n.service.Nodes[c.RemoteAddr().String()] = content
			log.Infof("Currently there are %d connected nodes", len(n.service.Nodes))
		}

		// When the node emit a ping message, we need to respond with pong
		// before five seconds to authorize that node to sent reports
		if msgType == messagePing {
			ping, err := parseNodePingMessage(msg)
			if err != nil {
				log.Warningf("Can't parse ping message sent by node[%s], error: %s", ping.ID, err)
				return
			}
			sendError := ping.SendResponse(c)
			if sendError != nil {
				log.Errorf("Error sending pong response to node[%s], error: %s", ping.ID, sendError)
			}
			n.service.Message <- content
		}

		// Send the content sent by the nodes directly to the consumer clients.
		// Only message types recognized by this server
		if isValidMessage(msgType) {
			n.service.Message <- content
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
