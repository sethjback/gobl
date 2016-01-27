package manager

import (
	"errors"
	"strings"

	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/spec"
)

// Agents returns a list of all known agents
func Agents() ([]*spec.Agent, error) {
	return gDb.AgentList()
}

// AddAgent simply adds an agent to the system
func AddAgent(agent *spec.Agent) error {

	key, err := getAgentKey(agent)
	if err != nil {
		return err
	}

	agent.PublicKey = key

	ks, err := keys.DecodePublicKeyString(key)
	if err != nil {
		return err
	}

	ip := strings.Split(agent.Address, ":")
	keyManager.PublicKeys[ip[0]] = ks

	return gDb.AddAgent(agent)
}

// GetAgent returns the stored agent information
func GetAgent(agentID int) (*spec.Agent, error) {
	return gDb.GetAgent(agentID)
}

// agentKey contacts the given agentID and retrieves it's public key
func getAgentKey(agent *spec.Agent) (string, error) {

	aR := &httpapi.APIRequest{Address: agent.Address}

	response, err := aR.GET("/key")
	if err != nil {
		return "", err
	}

	key, ok := response.Data["keyString"].(string)
	if !ok {
		return "", errors.New("Invalid key returned from agent")
	}

	return key, nil
}

// UpdateAgent updates the information for given agent id
func UpdateAgent(agent *spec.Agent) error {
	return nil
}
