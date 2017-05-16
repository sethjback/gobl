package manager

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/email"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/model"
)

// JobStatus reads the status from the DB
func JobStatus(id string) (*model.JobMeta, error) {
	var jobMeta *model.JobMeta
	var rerr error

	j, err := gDb.GetJob(id)
	if err != nil {
		return nil, err
	}

	if j.Meta.State == model.StateRunning || j.Meta.State == model.StateNotification {
		aR := &httpapi.Request{Host: j.Agent.Address, Path: "/jobs/" + id, Method: "GET"}
		response, err := aR.Send(signer)
		if err != nil {
			return nil, err
		}
		//todo: update status on our end if the job isn't found on the agent. Likely causes are the agent Shutdown
		// uncleanly and wasn't able to persist the job
		if response.Error != nil {
			return nil, response.Error
		}

		if jd, ok := response.Data[id]; ok {
			jobMeta = jd.(*model.JobMeta)
		} else {
			rerr = errors.New("Unable to find job on agent")
		}
	}

	return jobMeta, rerr
}

// JobList builds a list of jobs based on the given paramiters
func JobList(filters map[string]string) ([]model.Job, error) {
	list, err := gDb.GetJobs(filters)
	if list == nil {
		list = make([]model.Job, 0)
	}

	return list, err
}

// AddJobFile adds a file entry to the job
func AddJobFile(jobID string, jobFile model.JobFile) error {
	job, err := gDb.GetJob(jobID)
	if err != nil {
		return err
	}
	if job.Meta.State != model.StateRunning {
		return errors.New("Cannot add files to completed job")
	}

	return gDb.SaveJobFile(jobID, jobFile)
}

// FinishJob updates the job status in the DB and begins the file indexing process
func FinishJob(id string) error {
	job, err := gDb.GetJob(id)
	if err != nil {
		return err
	}

	job.Meta.State = model.StateFinished
	job.Meta.End = time.Now()

	gDb.SaveJob(*job)

	// Todo: index table for files lookup

	if conf.Email.Configured() {
		body := "Job Complete: " + job.ID + "\n"
		body += "Agent: " + job.Agent.Name + "\n"
		body += "Start: " + job.Meta.Start.String() + "\nEnd: " + job.Meta.End.String() + "\nDuration: " + fmt.Sprintf("%v", job.Meta.End.Sub(job.Meta.Start)) + "\n"
		body += "Message: " + job.Meta.Message + "\n\n"
		body += "Job Definition: " + fmt.Sprintf("%+v", job.Definition)

		email.SendEmail(conf.Email, body, "Job Report: "+job.ID)
	}

	return nil
}

// JobFiles pulls a list of files in job
func JobFiles(jobID string, filters map[string]string) ([]model.JobFile, error) {
	_, err := gDb.GetJob(jobID)
	if err != nil {
		return nil, err
	}

	return gDb.JobFiles(jobID, filters)
}

func JobDirectories(jobID, parent string) ([]string, error) {
	_, err := gDb.GetJob(jobID)
	if err != nil {
		return nil, err
	}

	if parent == "" {
		parent = "/"
	}

	return gDb.JobDirectories(jobID, parent)
}

func NewJob(jobRequest model.Job) (string, error) {
	jobRequest.ID = uuid.New().String()
	jobRequest.Meta = &model.JobMeta{State: model.StateNew, Start: time.Now().UTC()}

	err := gDb.SaveJob(jobRequest)
	if err != nil {
		return "", err
	}

	aR := &httpapi.Request{Host: jobRequest.Agent.Address, Path: "/jobs", Method: "POST"}
	jobRequest.Agent = nil
	err = aR.SetBody(jobRequest)
	if err != nil {
		return "", err
	}

	response, err := aR.Send(signer)
	if err != nil {
		return "", err
	}

	return jobRequest.ID, response.Error
}

func GetJob(jobID string) (*model.Job, error) {
	return gDb.GetJob(jobID)
}
