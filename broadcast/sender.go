package broadcast

import "net/http"

const Root string = "/"

// InfoSender is the responsible to send node state to registered clients
type InfoSender struct {
}

// HandleRequest handle all request from clients that are not Ethereum nodes
func (i *InfoSender) HandleRequest(w http.ResponseWriter, r *http.Request) {

}
