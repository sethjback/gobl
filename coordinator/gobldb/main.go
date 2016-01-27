package gobldb

import (
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/spec"
)

// QueryDateFormat defines the format we want dates in for the joblist queries
const QueryDateFormat = "2006-01-02 15:04"

// Database is the interface that must be implemented by the DB driver
type Database interface {
	Init(map[string]interface{}) error
	// AGENTS
	AddAgent(*spec.Agent) error
	AgentList() ([]*spec.Agent, error)
	GetAgent(int) (*spec.Agent, error)
	UpdateAgent(*spec.Agent) error

	// BACKUP DEFINITIONS
	AddBackupDefinition(*spec.BackupDefinition) error
	GetBackupDefinition(int) (*spec.BackupDefinition, error)
	DeleteBackupDefinition(int) error
	UpdateBackupDefinition(*spec.BackupDefinition) error
	BackupDefinitionList() ([]*spec.BackupDefinition, error)

	// MANUALLY RUN JOBS
	CreateBackupJob(*spec.BackupDefinition) (*spec.Job, error)
	CreateRestoreJob(int, *spec.RestoreRequest) (*spec.Job, error)

	// JOBS
	GetJob(int) (*spec.Job, error)
	AddJobFile(int, *spec.JobFile) error
	UpdateJob(*spec.Job) error
	JobErrorCount(int) (int, error)
	JobQuery(map[string]string) ([]*spec.Job, error)
	JobFiles(string, int) ([]*spec.JobFile, error)
	JobFileSignatures(int, []int) ([]files.Signature, error)

	// SCHEDULES
	AddSchedule(*spec.Schedule) error
	DeleteSchedule(int) error
	UpdateSchedule(*spec.Schedule) error
	ScheduleList() ([]*spec.Schedule, error)
	GetSchedule(int) (*spec.Schedule, error)
}
