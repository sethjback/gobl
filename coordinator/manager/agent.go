package manager

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/certificates"
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
	if agent.Key == nil {
		caKey, _ := gDb.GetKey("CA")
		if caKey == nil {
			return "", errors.New("must create or set CA Key first")
		}
		key, err := certificates.NewHostCertificate(*caKey, agent.Name)
		if err != nil {
			return "", err
		}
		agent.Key = key
	}

	agent.ID = uuid.New().String()

	err := gDb.SaveAgent(agent)
	if err != nil {
		return "", err
	}

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

	return getAgentStatus(*agent)
}

func getAgentStatus(agent model.Agent) (map[string]interface{}, error) {
	return nil, nil
}

// UpdateAgent updates the information for given agent id
func UpdateAgent(agent model.Agent) error {
	current, err := gDb.GetAgent(agent.ID)
	if err != nil {
		return err
	}

	if agent.Key == nil {
		if current.Name != agent.Name {
			caKey, _ := gDb.GetKey("CA")
			if caKey == nil {
				return errors.New("must create or set CA Key first")
			}
			key, err := certificates.NewHostCertificate(*caKey, agent.Name)
			if err != nil {
				return err
			}
			agent.Key = key
		} else {
			agent.Key = current.Key
		}
	}

	return gDb.SaveAgent(agent)
}
