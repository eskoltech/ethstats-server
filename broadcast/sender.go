package broadcast

import (
	"log"
	"net/http"
	"strings"

	"github.com/eskoltech/ethstats/service"
	"github.com/gorilla/websocket"
)

// Root is the home endpoint where clients are registered to receive node updates
const Root string = "/"

// upgradeConnection upgrade only HTTP request to the / endpoint
var upgradeConnection = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return strings.Compare(r.RequestURI, Root) == 0
	},
}

// InfoSender is the responsible to send node state to registered clients
type InfoSender struct {
	clients map[*websocket.Conn]bool
	service *service.Channel
}

// New creates a new InfoSender struct with the required service
func New(service *service.Channel) *InfoSender {
	return &InfoSender{
		clients: make(map[*websocket.Conn]bool),
		service: service,
	}
}

// HandleRequest handle all request from clients that are not Ethereum nodes
func (i *InfoSender) HandleRequest(w http.ResponseWriter, r *http.Request) {
	clientConn, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
	}
	i.clients[clientConn] = true
	go i.loop()
}

// loop loops as the connection with the client is alive
func (i *InfoSender) loop() {
	for {
		msg := <-i.service.Message
		if len(i.clients) == 0 {
			break
		}
		i.writeMessage(msg)
	}
}

func (i *InfoSender) writeMessage(msg []byte) {
	for client := range i.clients {
		err := client.WriteMessage(1, msg)
		if err != nil {
			// close and delete the client connection and release
			client.Close()
			delete(i.clients, client)
		}
	}
}