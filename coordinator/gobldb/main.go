package gobldb

import (
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/model"
)

// QueryDateFormat defines the format we want dates in for the joblist queries
const QueryDateFormat = "2006-01-02 15:04"

// Database is the interface that must be implemented by the DB driver
type Database interface {
	Init(config.DB) error
	Close() error

	// AGENTS
	SaveAgent(*model.Agent) error
	GetAgent(id string) (*model.Agent, error)
	AgentList() ([]model.Agent, error)

	// Job Definitions
	SaveJobDefinition(model.JobDefinition) error
	DeleteJobDefinition(model.JobDefinition) error
	GetJobDefinition(string) (*model.JobDefinition, error)
	GetJobDefinitions() ([]model.JobDefinition, error)

	// JOBS
	SaveJob(job model.Job) error
	GetJob(id string) (*model.Job, error)
	GetJobs(filters map[string]string) ([]model.Job, error)

	AddJobFile(jobID string, jobfile model.JobFile) error
	JobFiles(filters map[string]string) ([]model.JobFile, error)
	JobDirectories(jobID, parent string) ([]string, error)

	// SCHEDULES
	SaveSchedule(*model.Schedule) error
	DeleteSchedule(id string) error
	ScheduleList() ([]model.Schedule, error)
}
