package entities

// Conn defines the basic structure of a connection
// from a node to another node
type Conn struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}
