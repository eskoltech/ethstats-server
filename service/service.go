package service

// Channel is the service whereby servers exchange info
type Channel struct {
	Message chan []byte
}
