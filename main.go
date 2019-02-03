package main

import (
	"flag"
	"fmt"
	"github.com/eskoltech/ethstats/message"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

const (
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
		msg := message.Emit{Content: content}
		rsp, err := msg.Response()
		if err != nil {
			log.Print(err)
			break
		}
		err = c.WriteMessage(mt, rsp)
		if err != nil {
			log.Print(err)
			break
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
