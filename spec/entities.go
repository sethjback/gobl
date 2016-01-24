package spec

// Coordinator defines a coordinator
type Coordinator struct {
	Address   string `json:"address",toml:"address"`
	PublicKey string `json:"publickey"`
}

// Agent defines and agent
type Agent struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	PublicKey string `json:"publickey"`
}
