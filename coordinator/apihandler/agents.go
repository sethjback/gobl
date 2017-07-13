package apihandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/certificates"
	dberr "github.com/sethjback/gobl/gobldb/errors"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/model"
)

type UpdateAgent struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func agentList(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	list, err := db.AgentList()
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	if list == nil {
		list = make([]model.Agent, 0)
	}

	return httpapi.Response{Data: map[string]interface{}{"agents": list}, HTTPCode: 200}
}

func addAgent(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	var a model.Agent
	err := r.JsonBody(&a)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	if a.Name == "" {
		return httpapi.Response{Error: errors.New("agent name must not be empty"), HTTPCode: 400}
	}
	if a.Address == "" {
		return httpapi.Response{Error: errors.New("agent name must not be empty"), HTTPCode: 400}
		//TODO: IP validation
	}
	if a.Key == nil {
		caKey, _ := db.GetKey("CA")
		if caKey == nil {
			return httpapi.Response{Error: errors.New("must create or set CA Key first"), HTTPCode: 400}
		}
		key, kerr := certificates.NewHostCertificate(*caKey, a.Name)
		if kerr != nil {
			return httpapi.Response{Error: kerr, HTTPCode: 400}
		}
		a.Key = key
	}

	a.ID = uuid.New().String()

	err = db.SaveAgent(a)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	//TODO: add location header
	return httpapi.Response{Data: map[string]interface{}{"id": a.ID}, HTTPCode: 201}
}

func getAgent(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id := ps.ByName("id")

	_, e := uuid.Parse(id)
	if e != nil {
		return httpapi.Response{Error: e, HTTPCode: 400}
	}

	agent, err := db.GetAgent(id)
	if err != nil {
		gerr := err.(goblerr.Error)
		if gerr.Code == dberr.ErrCodeNotFound {
			return httpapi.Response{Error: gerr, HTTPCode: 404}
		}
		return httpapi.Response{Error: gerr, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"agent": agent}, HTTPCode: 200}
}

func agentStatus(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	// TODO: Implement
	return httpapi.Response{Data: map[string]interface{}{}, HTTPCode: 200}
}

func updateAgent(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id := ps.ByName("id")

	_, e := uuid.Parse(id)
	if e != nil {
		return httpapi.Response{Error: e, HTTPCode: 400}
	}

	var a UpdateAgent
	err := r.JsonBody(&a)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	current, err := db.GetAgent(id)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	update := model.Agent{}

	if a.Name != "" && current.Name != a.Name {
		caKey, _ := db.GetKey("CA")
		if caKey == nil {
			return httpapi.Response{Error: errors.New("must create or set CA Key first"), HTTPCode: 400}
		}
		key, kerr := certificates.NewHostCertificate(*caKey, a.Name)
		if kerr != nil {
			return httpapi.Response{Error: kerr, HTTPCode: 400}
		}
		update.Key = key
		update.Name = a.Name
	} else {
		update.Key = current.Key
		update.Name = current.Name
	}

	if a.Address != "" {
		update.Address = a.Address
	}

	update.ID = id

	err = db.SaveAgent(update)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 200}
}
