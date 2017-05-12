package apihandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/coordinator/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/model"
)

func scheduleList(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	list, err := manager.ScheduleList()
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	if list == nil {
		list = make([]model.Schedule, 0)
	}

	return httpapi.Response{Data: map[string]interface{}{"stored": list, "active": manager.CronSchedules()}, HTTPCode: 200}
}

func addSchedule(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	var sched model.Schedule

	err := r.JsonBody(&sched)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	sID, gerr := manager.NewSchedule(sched)
	if gerr != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"id": sID}, HTTPCode: 201}
}

func updateSchedule(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("jobID"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid jobID"), HTTPCode: 400}
	}

	var sched model.Schedule
	err = r.JsonBody(&sched)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	sched.ID = id.String()

	if err := manager.UpdateSchedule(sched); err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 200}
}

func deleteSchedule(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("jobID"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid jobID"), HTTPCode: 400}
	}

	err = manager.DeleteSchedule(id.String())
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 200}
}
