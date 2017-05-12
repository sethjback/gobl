package manager

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/model"
)

// Agents returns a list of all known agents
func GetAgents() ([]model.Agent, error) {
	aS, err := gDb.AgentList()
	if aS == nil {
		aS = make([]model.Agent, 0)
	}
	return aS, err
}

// AddAgent simply adds an agent to the system
func AddAgent(agent model.Agent) (string, error) {
	key, err := getAgentKey(agent, signer)
	if err != nil {
		return "", err
	}

	ks, err := keys.DecodePublicKeyString(key)
	if err != nil {
		return "", err
	}
	agent.PublicKey = key
	agent.ID = uuid.New().String()

	err = gDb.SaveAgent(agent)
	if err != nil {
		return "", err
	}

	verifiers[agent.ID] = keys.NewVerifier(ks)
	return agent.ID, nil
}

// GetAgent returns the stored agent information
func GetAgent(agentID string) (*model.Agent, error) {
	return gDb.GetAgent(agentID)
}

func GetAgentStatus(agentID string) (map[string]interface{}, error) {
	agent, err := gDb.GetAgent(agentID)
	if err != nil {
		return nil, err
	}

	return getAgentStatus(*agent, signer)
}

func getAgentStatus(agent model.Agent, s keys.Signer) (map[string]interface{}, error) {
	aR := httpapi.NewRequest(agent.Address, "/status", "GET")

	response, err := aR.Send(s)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// agentKey contacts the given agentID and retrieves it's public key
func getAgentKey(agent model.Agent, s keys.Signer) (string, error) {
	aR := httpapi.NewRequest(agent.Address, "/key", "GET")
	//aR := &httpapi.Request{Host: agent.Address, Path: "/key", Method: "GET"}

	response, err := aR.Send(s)
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
func UpdateAgent(agent model.Agent, loadKey bool) error {
	current, err := gDb.GetAgent(agent.ID)
	if err != nil {
		return err
	}

	if loadKey {
		key, err := getAgentKey(agent, signer)
		if err != nil {
			return err
		}

		ks, err := keys.DecodePublicKeyString(key)
		if err != nil {
			return err
		}
		agent.PublicKey = key
		verifiers[agent.ID] = keys.NewVerifier(ks)
	} else {
		agent.PublicKey = current.PublicKey
	}

	return gDb.SaveAgent(agent)
}
