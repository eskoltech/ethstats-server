package broadcast

import (
	"github.com/eskoltech/ethstats/service"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

// hub maintain a list of registered clients to send messages
type hub struct {
	register chan *websocket.Conn
	close    chan interface{}
	clients  map[*websocket.Conn]bool
	service  *service.Channel
}

// loop loops as the server is alive and send messages to registered clients
func (h *hub) loop() {
	nodesReport := time.NewTicker(15 * time.Second)
	for {
		select {
		case msg := <-h.service.Message:
			if len(h.clients) == 0 {
				continue
			}
			h.writeMessage(msg)
		case client := <-h.register:
			h.clients[client] = true
		case <-h.close:
			h.quit()
			break
		case <-nodesReport.C:
			if len(h.service.Nodes) == 0 {
				// Don't send anything if no nodes connected...
				continue
			}
			for _, v := range h.service.Nodes {
				h.writeMessage(v)
			}
		}
	}
}

// writeMessage to all registered clients. If an error occurs sending a message to a client,
// then these connection is closed and removed from the pool of registered clients
func (h *hub) writeMessage(msg []byte) {
	for client := range h.clients {
		err := client.WriteMessage(1, msg)
		if err != nil {
			log.Infof("Closed connection with client: %s", client.RemoteAddr())
			// close and delete the client connection and release
			client.Close()
			delete(h.clients, client)
		}
	}
}

func (h *hub) quit() {
	log.Info("Closing all registered clients")
	for client := range h.clients {
		client.Close()
		delete(h.clients, client)
	}
	close(h.register)
	close(h.close)
}
