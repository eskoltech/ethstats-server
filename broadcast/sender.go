package broadcast

import (
	"net/http"
	"strings"

	"github.com/eskoltech/ethstats/service"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// Root is the home endpoint where clients are registered to receive node updates
const Root string = "/"

// upgradeConnection upgrade only HTTP request to the / endpoint
var upgradeConnection = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return strings.Compare(r.RequestURI, Root) == 0
	},
}

// Server is the responsible to send node state to registered clients
type Server struct {
	clients map[*websocket.Conn]bool
	service *service.Channel
}

// New creates a new Server struct with the required service
func New(service *service.Channel) *Server {
	return &Server{
		clients: make(map[*websocket.Conn]bool),
		service: service,
	}
}

// HandleRequest handle all request from clients that are not Ethereum nodes
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	clientConn, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error trying to establish communication with client (addr=%s, host=%s, URI=%s), %s",
			r.RemoteAddr, r.Host, r.RequestURI, err)
		return
	}
	s.clients[clientConn] = true
	log.Infof("Connected new client! (host=%s)", r.Host)
	go s.loop()
}

// loop loops as the connection with the client is alive
func (s *Server) loop() {
	for {
		msg := <-s.service.Message
		if len(s.clients) == 0 {
			log.Warning("No clients available to send stats from nodes")
			break
		}
		s.writeMessage(msg)
	}
}

func (s *Server) writeMessage(msg []byte) {
	for client := range s.clients {
		err := client.WriteMessage(1, msg)
		if err != nil {
			// close and delete the client connection and release
			client.Close()
			delete(s.clients, client)
		}
	}
}
