package apihandler

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sethjback/gobble/coordinator/manager"
	"github.com/sethjback/gobble/httpapi"
	"github.com/sethjback/gobble/spec"
)

func agentList(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	list, err := manager.Agents()
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not Read Agents", 500)
	}
	if list == nil {
		list = make([]*spec.Agent, 0)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"agents": list}, HTTPCode: 200}, nil
}

func addAgent(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	agent := new(spec.Agent)
	if err := json.Unmarshal(body, &agent); err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid Request", 400)
	}

	err = manager.AddAgent(agent)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Unable to add agent", 400)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"id": agent.ID}, HTTPCode: 200}, nil
}

func getAgent(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	vars := mux.Vars(r)

	aID, err := strconv.ParseInt(vars["agentID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid Agent ID", 400)
	}

	agent, err := manager.GetAgent(int(aID))
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Unable to get agent", 400)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"agent": agent}, HTTPCode: 200}, nil
}

func updateAgent(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	agent := new(spec.Agent)
	if err := json.Unmarshal(body, &agent); err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid Request", 400)
	}

	vars := mux.Vars(r)

	aID, err := strconv.ParseInt(vars["agentID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid Agent ID", 400)
	}

	agent.ID = int(aID)

	if err := manager.UpdateAgent(agent); err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not update agent", 400)
	}

	return &httpapi.APIResponse{HTTPCode: 200}, nil
}
