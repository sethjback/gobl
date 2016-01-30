package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sethjback/gobl/engines"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/spec"
)

func RunRestore(rr *spec.RestoreRequest) (int, error) {
	if _, err := engines.GetRestoreEngine(rr.To); err != nil {
		return -1, err
	}

	backupJob, err := gDb.GetJob(rr.JobID)
	if err != nil {
		return -1, err
	}

	found := false
	fromEngine := strings.ToLower(rr.From.Name)
	for _, engine := range backupJob.Definition.(spec.BackupParamiter).Engines {
		fmt.Println(fromEngine, engine)
		if fromEngine == strings.ToLower(engine.Name) {
			found = true
			break
		}
	}

	if !found {
		return -1, errors.New("From engine must have been used in the original backup job")
	}

	agent, err := gDb.GetAgent(backupJob.AgentID)
	if err != nil {
		return -1, err
	}

	jobfiles, err := gDb.JobFileSignatures(rr.JobID, rr.Files)
	if err != nil {
		return -1, err
	}

	fmt.Println("JobFiles: ", jobfiles)

	job, err := gDb.CreateRestoreJob(backupJob.AgentID, rr)

	rp := &spec.RestoreParamiter{
		To:               rr.To,
		From:             rr.From,
		FileSignatures:   jobfiles,
		BackupParamiters: backupJob.Definition.(spec.BackupParamiter)}

	request := &spec.RestoreJobRequest{
		ID:         job.ID,
		Paramiters: *rp}

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

	_, err = req.POST("/restores")
	if err != nil {
		return -1, err
	}

	return job.ID, nil
}
