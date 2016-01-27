package manager

import (
	"encoding/json"

	"github.com/sethjback/gobl/engines"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/modifications"
	"github.com/sethjback/gobl/spec"
)

// AddBackup adds a backupdefinition to the system
func AddBackup(b *spec.BackupDefinition) error {
	if _, err := modifications.GetModifications(b.Paramiters.Modifications, false); err != nil {
		return err
	}
	if _, err := engines.GetBackupEngines(b.Paramiters.Engines); err != nil {
		return err
	}
	return gDb.AddBackupDefinition(b)
}

// ModifyBackup updates the given backupID
func ModifyBackup(b *spec.BackupDefinition) error {
	return gDb.UpdateBackupDefinition(b)
}

// GetBackup returns the particular backup definition
func GetBackup(backupID int) (*spec.BackupDefinition, error) {
	return gDb.GetBackupDefinition(backupID)
}

// Backups returns a list of the defined backup jobs
func Backups() ([]*spec.BackupDefinition, error) {
	return gDb.BackupDefinitionList()
}

// DeleteBackup removes bacukp definition from the DB
func DeleteBackup(backupID int) error {
	return gDb.DeleteBackupDefinition(backupID)
}

// RunBackup will run the given backupID immediately
func RunBackup(backupID int) (int, error) {
	b, err := gDb.GetBackupDefinition(backupID)
	if err != nil {
		return -1, err
	}

	agent, err := gDb.GetAgent(b.AgentID)
	if err != nil {
		return -1, err
	}

	jSpec, err := gDb.CreateBackupJob(b)
	if err != nil {
		return -1, err
	}

	request := &spec.BackupJobRequest{Coordinator: &spec.Coordinator{Address: hostConfig["IP"].(string) + ":" + hostConfig["PORT"].(string)}}
	request.ID = jSpec.ID
	request.Paramiters = b.Paramiters

	bString, err := json.Marshal(request)
	if err != nil {
		return -1, err
	}

	sig, err := keyManager.Sign(string(bString))
	if err != nil {
		return -1, err
	}

	req := &httpapi.APIRequest{
		Address:   agent.Address,
		Body:      bString,
		Signature: sig}

	_, err = req.POST("/backups")
	if err != nil {
		return -1, err
	}

	return jSpec.ID, nil
}
