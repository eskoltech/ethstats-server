package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/eskoltech/ethstats/message"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	MessageHello string = "hello"

	API     string = "/api"
	VERSION string = "v0.1.0"
	BANNER  string = `
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
		return strings.Compare(r.RequestURI, API) == 0
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

	// Server loop
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
		if msgType == MessageHello {
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
			// Get value from JSON to store it and process it later to calculate
			// node latency etc
			authMsg, err := parseAuthMessage(msg)
			if err != nil {
				log.Print(err)
				return
			}
			nodeInfo[authMsg.ID] = append(nodeInfo[authMsg.ID], authMsg.Info)
		}
	}
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
	fmt.Printf(BANNER, VERSION)

	// check if server secret is valid
	if *secret == "" {
		log.Fatalln("Server secret can't be empty")
	}
	log.Printf("Starting websocket server in %s", *addr)
	http.HandleFunc(API, handleRequest)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
