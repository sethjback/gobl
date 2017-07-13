package gobldb

import (
	"errors"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/gobldb/leveldb"
	"github.com/sethjback/gobl/model"
)

type key int

const Config key = 0

// Database is the interface that must be implemented by the DB driver
type Database interface {
	Close() error
	SaveConfig(config.Store, map[string]string) error
	Init(config.Store) error

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

	// Coordinator key
	GetKey(string) (*model.Key, error)
	SaveKey(string, model.Key) error
}

type dbConfig struct {
	driver string
}

func SaveConfig(cs config.Store, env map[string]string) error {
	dbc := &dbConfig{}
	for k, v := range env {
		switch k {
		case "DB_DRIVER":
			dbc.driver = v
		}
	}

	cs.Add(Config, dbc)

	var err error

	switch dbc.driver {
	case "leveldb":
		l := &leveldb.Leveldb{}
		err = l.SaveConfig(cs, env)
	default:
		err = errors.New("Invalid DB driver specified")
	}

	return err
}

func configFromStore(cs config.Store) *dbConfig {
	if dbc, ok := cs.Get(Config); ok {
		return dbc.(*dbConfig)
	}
	return nil
}

func Get(cs config.Store) (Database, error) {
	dbc := configFromStore(cs)
	if dbc == nil {
		return nil, errors.New("Unable to find database config")
	}

	switch dbc.driver {
	case "leveldb":
		l := &leveldb.Leveldb{}
		err := l.Init(cs)
		return l, err
	default:
		return nil, errors.New("Invalid DB driver specified")
	}
}
