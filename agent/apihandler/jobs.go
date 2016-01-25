package apihandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sethjback/gobl/agent/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/spec"
)

func jobList(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	jobs := manager.Status()

	return &httpapi.APIResponse{Data: jobs, HTTPCode: 200}, nil
}

func newRestoreJob(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	restoreRequest := new(spec.RestoreJobRequest)
	if err := json.Unmarshal(body, &restoreRequest); err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid Request", 400)
	}

	err = manager.NewRestore(*restoreRequest)
	if err != nil {
		if e, ok := err.(*manager.Error); ok {
			return nil, httpapi.NewError(e.Source, e.Message, e.Code)
		}
		return nil, errors.New("")
	}

	return &httpapi.APIResponse{HTTPCode: 201}, nil
}

func newBackupJob(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	backupRequest := new(spec.BackupJobRequest)
	if err := json.Unmarshal(body, &backupRequest); err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid Request", 400)
	}

	if err = manager.VerifySignature(backupRequest.Coordinator, body, r.Header.Get("x-gobl-signature")); err != nil {
		fmt.Println(err)
		return nil, httpapi.NewError(err.Error(), "Invalid Signature", 401)
	}

	err = manager.NewBackup(*backupRequest)

	if err != nil {
		if e, ok := err.(*manager.Error); ok {
			return nil, httpapi.NewError(e.Source, e.Message, e.Code)
		}
		return nil, errors.New("")
	}

	return &httpapi.APIResponse{HTTPCode: 201}, nil
}

func cancelJob(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["jobId"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid job ID", 400)
	}

	if _, err := manager.Cancel(int(id)); err != nil {
		return nil, httpapi.NewError(err.Error(), "Problem Canceling Job", 500)
	}

	return &httpapi.APIResponse{HTTPCode: 200}, nil
}

func jobStatus(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["jobId"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid job ID", 400)
	}

	jobStatus, found := manager.JobStatus(int(id))
	if !found {
		return nil, httpapi.NewError("", "I can't find that job ID", 404)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{vars["jobId"]: jobStatus}, HTTPCode: 200}, nil
}
