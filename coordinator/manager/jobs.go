package manager

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/sethjback/gobl/email"
	"github.com/sethjback/gobl/spec"
)

// JobStatus reads the status from the DB
func JobStatus(id int) (*spec.Job, error) {
	return gDb.GetJob(id)
}

// QueryJobs builds a list of jobs based on the given paramiters
func QueryJobs(params url.Values) ([]*spec.Job, error) {

	validP := make(map[string]string)

	for _, key := range []string{"agentid", "start", "end", "state", "limit", "offset"} {
		val := params.Get(key)
		if len(val) != 0 {
			validP[key] = val
		}
	}

	// limit and offset should always be set
	if _, ok := validP["limit"]; !ok {
		validP["limit"] = "10"
	}

	if _, ok := validP["offset"]; !ok {
		validP["offset"] = "0"
	}

	return gDb.JobQuery(validP)
}

// AddJobFile adds a file entry to the job
func AddJobFile(jobID int, fileRequest *spec.JobFileRequest) error {
	job, err := gDb.GetJob(jobID)
	if err != nil {
		return err
	}
	if job.State != spec.Running {
		return errors.New("Cannot add files to completed job")
	}
	return gDb.AddJobFile(jobID, &fileRequest.File)
}

// FinishJob updates the job status in the DB and begins the file indexing process
func FinishJob(id int) error {
	job, err := gDb.GetJob(id)
	if err != nil {
		return err
	}

	count, err := gDb.JobErrorCount(id)
	if err != nil {
		return err
	}

	if count > 0 {
		job.State = spec.Errors
		m := "Job finished with " + strconv.Itoa(count) + " error(s)"
		job.Message = &m
	} else {
		job.State = spec.Complete
	}

	t := time.Now()

	job.End = &t

	err = gDb.UpdateJob(job)
	if err != nil {
		return err
	}

	// Todo: index table for files lookup

	if conf.Email.Configured() {
		a, _ := gDb.GetAgent(job.AgentID)
		var msg string
		if job.Message != nil {
			msg = *job.Message
		} else {
			msg = ""
		}
		body := "Job Complete: " + strconv.Itoa(job.ID) + "\n"
		body += "Agent: " + a.Name + "\n"
		body += "Start: " + job.Start.String() + "\nEnd: " + job.End.String() + "\nDuration: " + fmt.Sprintf("%v", job.End.Sub(*job.Start)) + "\n"
		body += "Message: " + msg + "\n\n"
		body += "Job Definition: " + fmt.Sprintf("%+v", job.Definition)

		email.SendEmail(conf.Email, body, "Job Report: "+strconv.Itoa(job.ID))
	}

	return nil
}

// JobFiles pulls a list of files in job
func JobFiles(path string, jobID int) ([]*spec.JobFile, error) {
	return gDb.JobFiles(path, jobID)
}
