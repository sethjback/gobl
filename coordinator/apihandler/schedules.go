package apihandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/coordinator/cron"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/model"
)

func scheduleList(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	list, err := cron.ScheduleList()
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	if list == nil {
		list = make([]model.Schedule, 0)
	}

	return httpapi.Response{Data: map[string]interface{}{"stored": list, "active": cron.Active()}, HTTPCode: 200}
}

func addSchedule(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	var sched model.Schedule

	err := r.JsonBody(&sched)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	_, err = db.GetAgent(sched.AgentID)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	_, err = db.GetJobDefinition(sched.JobDefinitionID)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	sID, gerr := cron.NewSchedule(sched)
	if gerr != nil {
		return httpapi.Response{Error: gerr, HTTPCode: 400}
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

	_, err = db.GetAgent(sched.AgentID)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	_, err = db.GetJobDefinition(sched.JobDefinitionID)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	if err := cron.UpdateSchedule(sched); err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 200}
}

func deleteSchedule(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("jobID"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid jobID"), HTTPCode: 400}
	}

	err = cron.DeleteSchedule(id.String())
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 200}
}
