package main

import (
	"flag"
	"fmt"
	"github.com/eskoltech/ethstats/broadcast"
	"github.com/eskoltech/ethstats/relay"
	"log"
	"net/http"
)

const (
	version string = "v0.1.0\n"
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

	// Node relay configuration
	nodeRelay := &relay.NodeRelay{Secret: *secret}
	http.HandleFunc(relay.Api, nodeRelay.HandleRequest)

	// Info sender configuration
	infoSender := &broadcast.InfoSender{}
	http.HandleFunc(broadcast.Root, infoSender.HandleRequest)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
