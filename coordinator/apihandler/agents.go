package apihandler

import (
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/coordinator/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/model"
)

func agentList(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	list, err := manager.GetAgents()
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"agents": list}, HTTPCode: 200}
}

func addAgent(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	var a model.Agent
	err := r.JsonBody(&a)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	id, e := manager.AddAgent(a)
	if e != nil {
		return httpapi.Response{Error: e, HTTPCode: 400}
	}

	//TODO: add location header
	return httpapi.Response{Data: map[string]interface{}{"id": id}, HTTPCode: 201}
}

func getAgent(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id := ps.ByName("id")

	_, e := uuid.Parse(id)
	if e != nil {
		return httpapi.Response{Error: e, HTTPCode: 400}
	}

	agent, err := manager.GetAgent(id)
	if err != nil {
		//TODO: need to return 404 if not found
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"agent": agent}, HTTPCode: 200}
}

func agentStatus(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id := ps.ByName("id")

	_, e := uuid.Parse(id)
	if e != nil {
		return httpapi.Response{Error: e, HTTPCode: 400}
	}

	status, err := manager.GetAgentStatus(id)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"status": status}, HTTPCode: 200}
}

func updateAgent(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id := ps.ByName("id")

	_, e := uuid.Parse(id)
	if e != nil {
		return httpapi.Response{Error: e, HTTPCode: 400}
	}

	var a model.Agent
	err := r.JsonBody(&a)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}
	a.ID = id

	updateKey := false
	if uk := ps.ByName("updateKey"); uk != "" {
		updateKey = true
	}

	if err := manager.UpdateAgent(a, updateKey); err != nil {
		return httpapi.Response{Error: e, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 200}
}
