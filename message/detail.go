package message

// NodeInfo is the information related to the Ethereum node
type NodeInfo struct {
	Name     string `json:"name"`
	Node     string `json:"node"`
	Port     int    `json:"port"`
	Network  string `json:"net"`
	Protocol string `json:"protocol"`
	API      string `json:"api"`
	Os       string `json:"os"`
	OsVer    string `json:"os_v"`
	Client   string `json:"client"`
	History  bool   `json:"canUpdateHistory"`
}

// AuthMessage is the struct sent by the server on the first connection
type AuthMessage struct {
	ID     string   `json:"id"`
	Secret string   `json:"secret"`
	Info   NodeInfo `json:"info"`
}
