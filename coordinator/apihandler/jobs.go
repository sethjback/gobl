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

func jobQuery(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	jobs, err := manager.QueryJobs(r.URL.Query())
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Error with query", 400)
	}

	if jobs == nil {
		jobs = make([]*spec.Job, 0)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"jobs": jobs}, HTTPCode: 200}, nil
}

func jobStatus(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["jobID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid job id", 400)
	}
	job, err := manager.JobStatus(int(id))
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not get job status", 400)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"job": job}, HTTPCode: 200}, nil
}

func jobFiles(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["jobID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid job id", 400)
	}

	files, err := manager.JobFiles(vars["path"], int(id))
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not pull list of job files", 400)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"files": files}, HTTPCode: 200}, nil
}

func cancelJob(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	return nil, nil
}

func addJobFile(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	vars := mux.Vars(r)

	jID, err := strconv.ParseInt(vars["jobID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid job id", 400)
	}

	if err = manager.VerifySignature(r.RemoteAddr, body, r.Header.Get("x-gobl-signature")); err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid Signature", 401)
	}

	var jobFile spec.JobFileRequest
	if err = json.Unmarshal(body, &jobFile); err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid request", 400)
	}

	if err = manager.AddJobFile(int(jID), &jobFile); err != nil {
		return nil, httpapi.NewError(err.Error(), "Unable to add file to job", 400)
	}

	return &httpapi.APIResponse{HTTPCode: 200}, nil
}

func finishJob(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	if err = manager.VerifySignature(r.RemoteAddr, body, r.Header.Get("x-gobl-signature")); err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid Signature", 401)
	}

	vars := mux.Vars(r)

	jID, err := strconv.ParseInt(vars["jobID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid job id", 400)
	}

	var jobUpdate spec.JobUpdateRequest
	if err = json.Unmarshal(body, &jobUpdate); err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid request", 400)
	}

	err = manager.FinishJob(int(jID))
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Problem updating Job", 400)
	}

	return &httpapi.APIResponse{HTTPCode: 200}, nil
}
