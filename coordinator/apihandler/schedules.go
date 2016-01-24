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

func scheduleList(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	list, err := manager.ScheduleList()
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not Read Agents", 500)
	}
	if list == nil {
		list = make([]*spec.Schedule, 0)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"stored": list, "active": manager.CronSchedules()}, HTTPCode: 200}, nil
}

func addSchedule(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	s := new(spec.Schedule)
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid Request", 400)
	}

	err = manager.AddSchedule(s)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Unable to add schedule", 400)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"id": s.ID}, HTTPCode: 200}, nil
}

func getSchedule(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	vars := mux.Vars(r)

	sID, err := strconv.ParseInt(vars["sID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid schedule ID", 400)
	}

	s, err := manager.GetSchedule(int(sID))
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Unable to get schedule", 400)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"schedule": s}, HTTPCode: 200}, nil
}

func updateSchedule(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	s := new(spec.Schedule)
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid Request", 400)
	}

	vars := mux.Vars(r)

	sID, err := strconv.ParseInt(vars["sID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid schedule ID", 400)
	}

	s.ID = int(sID)

	if err := manager.UpdateSchedule(s); err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not update schedule", 400)
	}

	return &httpapi.APIResponse{HTTPCode: 200}, nil
}

func deleteSchedule(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	vars := mux.Vars(r)

	sID, err := strconv.ParseInt(vars["sID"], 10, 64)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid schedule ID", 400)
	}

	err = manager.DeleteSchedule(int(sID))
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Could not delete schedule", 500)
	}
	return &httpapi.APIResponse{HTTPCode: 200}, nil
}
