package manager

import (
	"github.com/google/uuid"
	"github.com/sethjback/gobl/model"
)

func CreateJobDefinition(jobDef model.JobDefinition) (string, error) {
	jobDef.ID = uuid.New().String()

	return jobDef.ID, gDb.SaveJobDefinition(jobDef)
}

func UpdateJobDefinition(jobDef model.JobDefinition) error {
	return gDb.SaveJobDefinition(jobDef)
}

func DeleteJobDefinition(id string) error {
	return gDb.DeleteJobDefinition(id)
}

func GetJobDefinition(id string) (*model.JobDefinition, error) {
	return gDb.GetJobDefinition(id)
}

func GetJobDefinitions() ([]model.JobDefinition, error) {
	jdefs, err := gDb.JobDefinitionList()
	if jdefs == nil {
		jdefs = make([]model.JobDefinition, 0)
	}
	return jdefs, err
}
