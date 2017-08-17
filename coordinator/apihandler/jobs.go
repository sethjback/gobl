package apihandler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/coordinator/grpcclient"
	pb "github.com/sethjback/gobl/goblgrpc"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/model"
)

type JobRequest struct {
	Definition model.JobDefinition `json:"jobDefinition"`
	Agent      string              `json:"agentId"`
}

func jobList(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	jobs, err := db.JobList(queryToMap(r.Query))
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

	j, err := db.GetJob(id.String())
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	if j.Meta.State == model.StateRunning || j.Meta.State == model.StateNotification {
		// TODO: implement via grpc
		//todo: update status on our end if the job isn't found on the agent. Likely causes are the agent Shutdown
		// uncleanly and wasn't able to persist the job
	}

	return httpapi.Response{Data: map[string]interface{}{"job": j}, HTTPCode: 200}
}

func jobFiles(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		return httpapi.Response{Error: errors.New("Invalid job id"), HTTPCode: 400}
	}

	_, err = db.GetJob(id.String())
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	files, err := db.JobFileList(id.String(), queryToMap(r.Query))
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

	_, err = db.GetJob(id.String())
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	parent := r.Query.Get("parent")
	if parent == "" {
		parent = "/"
	}

	dirs, err := db.JobDirectories(id.String(), parent)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"directories": dirs}, HTTPCode: 200}
}

func cancelJob(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	return httpapi.Response{Error: errors.New("Unimplemented"), HTTPCode: 400}
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

	agent, err := db.GetAgent(jr.Agent)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	aClient, err := grpcclient.New(agent)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	//TODO: validate jobRequest

	job := model.Job{
		ID:         uuid.New().String(),
		Meta:       &model.JobMeta{State: model.StateNew, Start: time.Now().UTC()},
		Agent:      agent,
		Definition: &jr.Definition,
	}

	err = db.SaveJob(job)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	pbdef := buildProtoJobDef(*job.Definition)
	pbdef.Id = job.ID
	var rmsg *pb.ReturnMessage
	if job.Definition.Type == "backup" {
		rmsg, err = aClient.Pb.Backup(context.Background(), pbdef)
	} else {
		rmsg, err = aClient.Pb.Restore(context.Background(), pbdef)
	}

	aClient.ClientConn.Close()

	if err != nil {
		job.Meta.Message = err.Error()
		job.Meta.State = model.StateFailed
		db.SaveJob(job)
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	if rmsg != nil && rmsg.Message != "success" {
		job.Meta.Message = rmsg.Message
		job.Meta.State = model.StateFailed
		db.SaveJob(job)
		return httpapi.Response{Error: errors.New(rmsg.String()), HTTPCode: 400}
	}

	return httpapi.Response{Data: map[string]interface{}{"id": job.ID}, HTTPCode: 201}
}

func buildProtoJobDef(jd model.JobDefinition) *pb.JobDefinition {
	pbdef := &pb.JobDefinition{}
	for _, v := range jd.To {
		t := &pb.MEDefinition{}
		t.Name = v.Name
		for n, v := range v.Options {
			t.Options = append(t.Options, &pb.Option{Name: n, Value: v})
		}

		pbdef.To = append(pbdef.To, t)
	}

	if jd.From != nil {
		t := &pb.MEDefinition{}
		t.Name = jd.From.Name
		for n, v := range jd.From.Options {
			t.Options = append(t.Options, &pb.Option{Name: n, Value: v})
		}

		pbdef.From = t
	}

	for _, v := range jd.Modifications {
		t := &pb.MEDefinition{}
		t.Name = v.Name
		for n, v := range v.Options {
			t.Options = append(t.Options, &pb.Option{Name: n, Value: v})
		}

		pbdef.Modifications = append(pbdef.Modifications, t)
	}

	for _, v := range jd.Paths {
		t := &pb.Path{}
		t.Root = v.Root
		t.Excludes = v.Excludes
		pbdef.Paths = append(pbdef.Paths, t)
	}

	return pbdef
}
