package apihandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/coordinator/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/model"
)

func createJobDefinition(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	var jd model.JobDefinition
	gerr := r.JsonBody(&jd)
	if gerr != nil {
		return httpapi.Response{Error: gerr, HTTPCode: 400}
	}

	id, err := manager.CreateJobDefinition(jd)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"id": id}, HTTPCode: 201}
}

func updateJobDefinition(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid job id"), HTTPCode: 400}
	}

	var jd model.JobDefinition
	gerr := r.JsonBody(&jd)
	if gerr != nil {
		return httpapi.Response{Error: gerr, HTTPCode: 400}
	}

	jd.ID = id.String()

	err = manager.UpdateJobDefinition(jd)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 200}
}

func deleteJobDefinition(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid job id"), HTTPCode: 400}
	}

	err = manager.DeleteJobDefinition(id.String())
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 200}
}

func getJobDefinition(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid job id"), HTTPCode: 400}
	}

	jd, err := manager.GetJobDefinition(id.String())
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"jobDefinition": jd}, HTTPCode: 200}
}

func jobDefinitionList(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	jdefs, err := manager.GetJobDefinitions()
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"jobDefinitions": jdefs}, HTTPCode: 200}
}
