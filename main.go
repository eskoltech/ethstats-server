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

var addr = flag.String("addr", "localhost:3000", "http service address")

// upgradeConnection allows
var upgradeConnection = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return strings.Compare(r.RequestURI, API) == 0
	},
}

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

		// If message type is hello, we need to check if the secret is
		// correct, and then, send a ready message
		if msg.GetType() == MessageHello {
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
		}
	}
}

func main() {
	flag.Parse()
	fmt.Printf(BANNER, VERSION)
	log.Printf("Starting websocket server in %s", *addr)

	http.HandleFunc(API, handleRequest)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
