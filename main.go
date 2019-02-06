package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/eskoltech/ethstats/message"
	"github.com/gorilla/websocket"
)

const (
	messageHello string = "hello"
	messagePing         = "node-ping"

	api     string = "/api"
	version string = "v0.1.0"
	banner  string = `
        __  .__              __          __          
  _____/  |_|  |__   _______/  |______ _/  |_  ______
_/ __ \   __\  |  \ /  ___/\   __\__  \\   __\/  ___/
\  ___/|  | |   Y  \\___ \  |  |  / __ \|  |  \___ \ 
 \___  >__| |___|  /____  > |__| (____  /__| /____  >
     \/          \/     \/            \/          \/  %s
`
)

var addr = flag.String("addr", "localhost:3000", "HTTP service address")
var secret = flag.String("secret", "", "Server secret")

// upgradeConnection allows
var upgradeConnection = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return strings.Compare(r.RequestURI, api) == 0
	},
}

// nodeInfo is the list of current connected nodes
var nodeInfo = make(map[string][]message.NodeInfo)

// handleRequest is the function to handle all server requests...
func handleRequest(w http.ResponseWriter, r *http.Request) {
	c, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Close connection if an unexpected error occurs
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(c)

	for {
		mt, content, err := c.ReadMessage()
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
			log.Printf("Auth message from '%s' node for network %s, node=%s",
				authMsg.ID, authMsg.Info.Network, authMsg.Info.Node)

			// first check if the secret is correct
			if authMsg.Secret != *secret {
				log.Printf("Invalid secret from node %s", authMsg.ID)
				return
			}
			ready := map[string][]interface{}{"emit": {"ready"}}
			response, err := json.Marshal(ready)
			if err != nil {
				log.Print(err)
				return
			}
			err = c.WriteMessage(mt, response)
			if err != nil {
				log.Print(err)
				return
			}
			// When the message is correctly sent, save the node
			nodeInfo[authMsg.ID] = append(nodeInfo[authMsg.ID], authMsg.Info)
		}

		// When the node emit a ping message, we need to respond with pong
		// before five seconds to authorize that node to sent reports
		if msgType == messagePing {
			ping, err := parseNodePingMessage(msg)
			if err != nil {
				log.Print(err)
				return
			}
			log.Printf("Received ping from '%s'", ping.ID)
			sendError := ping.SendResponse(c)
			if sendError != nil {
				log.Print(sendError)
			}
			log.Printf("Server sent pong to node '%s'", ping.ID)
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
func parseAuthMessage(msg message.Message) (message.AuthMessage, error) {
	value, err := msg.GetValue()
	if err != nil {
		return message.AuthMessage{}, err
	}
	var detail message.AuthMessage
	err = json.Unmarshal(value, &detail)
	return detail, err
}

// main is the program entry point. If the server secret is not set when
// init, the server can't start
func main() {
	flag.Parse()
	fmt.Printf(banner, version)

	// check if server secret is valid
	if *secret == "" {
		log.Fatalln("Server secret can't be empty")
	}
	log.Printf("Starting websocket server in %s", *addr)
	http.HandleFunc(api, handleRequest)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
