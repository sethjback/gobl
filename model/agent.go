package model

// Agent defines and agent
type Agent struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	PublicKey string `json:"publickey"`
}
