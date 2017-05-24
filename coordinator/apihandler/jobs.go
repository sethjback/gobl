package apihandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/coordinator/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/model"
)

type JobRequest struct {
	Definition model.JobDefinition `json:"jobDefinition"`
	Agent      string              `json:"agentId"`
}

func jobList(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	jobs, err := manager.JobList(queryToMap(r.Query))
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"jobs": jobs}, HTTPCode: 200}
}

func jobStatus(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid job id"), HTTPCode: 400}
	}

	job, err := manager.JobStatus(id.String())
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"job": job}, HTTPCode: 200}
}

func jobFiles(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid job id"), HTTPCode: 400}
	}

	files, err := manager.JobFiles(id.String(), queryToMap(r.Query))
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"files": files}, HTTPCode: 200}
}

func jobDirectories(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid job id"), HTTPCode: 400}
	}

	dirs, err := manager.JobDirectories(id.String(), r.Query.Get("parent"))
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"directories": dirs}, HTTPCode: 200}
}

func cancelJob(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	return httpapi.Response{Error: errors.New("Unimplemented"), HTTPCode: 400}
}

func addJobFile(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	var jf model.JobFile
	gerr := r.JsonBody(&jf)
	if gerr != nil {
		return httpapi.Response{Error: gerr, HTTPCode: 400}
	}

	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid job id"), HTTPCode: 400}
	}

	if err = manager.AddJobFile(id.String(), jf); err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 201}
}

func finishJob(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid job id"), HTTPCode: 400}
	}

	err = manager.FinishJob(id.String())
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 200}
}

func newJob(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	var jr JobRequest
	gerr := r.JsonBody(&jr)
	if gerr != nil {
		return httpapi.Response{Error: gerr, HTTPCode: 400}
	}

	aID, err := uuid.Parse(jr.Agent)
	if err != nil {
		return httpapi.Response{Error: errors.New("Unable to parse agent ID: " + err.Error()), HTTPCode: 400}
	}

	jr.Agent = aID.String()

	//TODO: validate jobRequest
	id, err := manager.NewJob(jr.Definition, jr.Agent)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"id": id}, HTTPCode: 201}
}
