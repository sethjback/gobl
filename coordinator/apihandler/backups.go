package apihandler

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sethjback/gobl/coordinator/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/spec"
)

func addBackup(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	backup := new(spec.BackupDefinition)

	err = json.Unmarshal(body, &backup)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid backup definition provided", 400)
	}

	err = manager.AddBackup(backup)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not create backup", 400)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"id": backup.ID}, HTTPCode: 201}, nil
}

func modifyBackup(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	backup := new(spec.BackupDefinition)

	err = json.Unmarshal(body, &backup)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid backup definition provided", 400)
	}

	vars := mux.Vars(r)

	bID, err := strconv.ParseInt(vars["backupID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid backup ID", 400)
	}

	backup.ID = int(bID)

	err = manager.ModifyBackup(backup)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Unable to modify backup", 400)
	}

	return &httpapi.APIResponse{HTTPCode: 200}, nil
}

func runBackup(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	vars := mux.Vars(r)

	bID, err := strconv.ParseInt(vars["backupID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid backup ID", 400)
	}

	id, err := manager.RunBackup(int(bID))
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Failed to run backup", 400)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"id": id}, HTTPCode: 200}, nil
}

func listBackups(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	list, err := manager.Backups()
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not read backups", 500)
	}
	if list == nil {
		list = make([]*spec.BackupDefinition, 0)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"backups": list}, HTTPCode: 200}, nil
}

func getBackup(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	vars := mux.Vars(r)

	bID, err := strconv.ParseInt(vars["backupID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid backup ID", 400)
	}

	backup, err := manager.GetBackup(int(bID))
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not read backup", 500)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"backup": backup}, HTTPCode: 200}, nil
}

func deleteBackup(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	vars := mux.Vars(r)

	bID, err := strconv.ParseInt(vars["backupID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid backup ID", 400)
	}

	err = manager.DeleteBackup(int(bID))
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not delete backup", 500)
	}
	return &httpapi.APIResponse{HTTPCode: 200}, nil
}
