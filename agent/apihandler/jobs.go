package apihandler

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/agent/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/model"
)

func jobList(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	jobs := manager.Status()

	return httpapi.Response{Data: jobs, HTTPCode: 200}
}

func newJob(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	var job model.Job

	err := r.JsonBody(&job)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	if job.Definition.Type == model.TypeBackup {
		err = manager.NewBackup(job)
	} else {
		err = manager.NewRestore(job)
	}

	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 201}
}

func cancelJob(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id := ps.ByName("id")

	err := manager.Cancel(id)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 200}
}

func jobStatus(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id := ps.ByName("id")

	jobStatus, err := manager.JobStatus(id)
	if err != nil {
		//TODO: need to return 404 if not found
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{id: jobStatus}, HTTPCode: 200}
}
