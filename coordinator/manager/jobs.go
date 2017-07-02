package manager

import (
	"errors"
	"time"

	"github.com/google/uuid"
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
		// TODO: implement via grpc
		//todo: update status on our end if the job isn't found on the agent. Likely causes are the agent Shutdown
		// uncleanly and wasn't able to persist the job

	}

	return jobMeta, rerr
}

// JobList builds a list of jobs based on the given paramiters
func JobList(filters map[string]string) ([]model.Job, error) {
	list, err := gDb.JobList(filters)
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
	job.Meta.End = time.Now().UTC()

	gDb.SaveJob(*job)

	/*
		if conf.Email.Configured() {
			body := "Job Complete: " + job.ID + "\n"
			body += "Agent: " + job.Agent.Name + "\n"
			body += "Start: " + job.Meta.Start.String() + "\nEnd: " + job.Meta.End.String() + "\nDuration: " + fmt.Sprintf("%v", job.Meta.End.Sub(job.Meta.Start)) + "\n"
			body += "Message: " + job.Meta.Message + "\n\n"
			body += "Job Definition: " + fmt.Sprintf("%+v", job.Definition)

			email.SendEmail(conf.Email, body, "Job Report: "+job.ID)
		}
	*/
	return nil
}

// JobFiles pulls a list of files in job
func JobFiles(jobID string, filters map[string]string) ([]model.JobFile, error) {
	_, err := gDb.GetJob(jobID)
	if err != nil {
		return nil, err
	}

	return gDb.JobFileList(jobID, filters)
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

func NewJob(jobDefinition model.JobDefinition, agentID string) (string, error) {
	agent, err := gDb.GetAgent(agentID)
	if err != nil {
		return "", err
	}

	job := model.Job{
		ID:         uuid.New().String(),
		Meta:       &model.JobMeta{State: model.StateNew, Start: time.Now().UTC()},
		Agent:      agent,
		Definition: &jobDefinition,
	}

	err = gDb.SaveJob(job)
	if err != nil {
		return "", err
	}

	//TODO: implement via grpc

	return job.ID, nil
}

func GetJob(jobID string) (*model.Job, error) {
	return gDb.GetJob(jobID)
}
