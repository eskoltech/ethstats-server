package broadcast

import (
	"net/http"
	"strings"

	"github.com/eskoltech/ethstats/service"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// Root is the home endpoint where hub are registered to receive node updates
const Root string = "/"

// upgradeConnection upgrade only HTTP request to the / endpoint
var upgradeConnection = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return strings.Compare(r.RequestURI, Root) == 0
	},
}

// Server is the responsible to send node state to registered hub
type Server struct {
	hub *hub
}

// New creates a new Server struct with the required service
func New(service *service.Channel) *Server {
	hub := &hub{
		register: make(chan *websocket.Conn),
		close:    make(chan interface{}),
		clients:  make(map[*websocket.Conn]bool),
		service:  service,
	}
	go hub.loop()
	return &Server{hub: hub}
}

// Close this server and all registered client connections
func (s *Server) Close() {
	s.hub.close <- "close"
}

// HandleRequest handle all request from hub that are not Ethereum nodes
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	clientConn, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error trying to establish communication with client (addr=%s, host=%s, URI=%s), %s",
			r.RemoteAddr, r.Host, r.RequestURI, err)
		return
	}
	log.Infof("Connected new client! (host=%s)", r.Host)
	s.hub.register <- clientConn
}
