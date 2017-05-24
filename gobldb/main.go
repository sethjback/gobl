package gobldb

import "github.com/sethjback/gobl/model"

// QueryDateFormat defines the format we want dates in
const QueryDateFormat = "2006-01-02 15:04"

// Database is the interface that must be implemented by the DB driver
type Database interface {
	Close() error

	// AGENTS
	SaveAgent(model.Agent) error
	GetAgent(id string) (*model.Agent, error)
	AgentList() ([]model.Agent, error)

	// Job Definitions
	SaveJobDefinition(model.JobDefinition) error
	DeleteJobDefinition(id string) error
	GetJobDefinition(string) (*model.JobDefinition, error)
	JobDefinitionList() ([]model.JobDefinition, error)

	// JOBS
	SaveJob(job model.Job) error
	GetJob(id string) (*model.Job, error)
	JobList(filters map[string]string) ([]model.Job, error)

	SaveJobFile(jobID string, jobfile model.JobFile) error
	JobFileList(jobID string, filters map[string]string) ([]model.JobFile, error)
	JobDirectories(jobID, parent string) ([]string, error)

	// SCHEDULES
	SaveSchedule(model.Schedule) error
	GetSchedule(id string) (*model.Schedule, error)
	DeleteSchedule(id string) error
	ScheduleList() ([]model.Schedule, error)

	// USERS
	GetUser(email string) (*model.User, error)
	UserList() ([]model.User, error)
	SaveUser(user model.User) error
	DeleteUser(email string) error
}
