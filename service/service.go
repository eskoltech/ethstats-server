package service

// Channel is the service whereby servers exchange info
type Channel struct {
	// Message is the content of the stats reported by the Ethereum node
	Message chan []byte
}
