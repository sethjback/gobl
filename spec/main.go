package spec

import (
	"fmt"
	"time"

	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/modification"
)

const (
	// Running state
	Running = 1
	// Errors encountered
	Errors = 2
	// Complete state
	Complete = 3
)

// JobFile records information about the files that are a part of a backup
type JobFile struct {
	ID        int             `json:"id,omitempty"`
	Signature files.Signature `json:"signature"`
	State     int             `json:"state"`
	Message   string          `json:"message,omitempty"`
}

// BackupJobRequest defines a request for a job sent to the AgentID
type BackupJobRequest struct {
	ID         string          `json:"id"`
	Paramiters BackupParamiter `json:"paramiters"`
}

// RestoreJobRequest defines a request for a job sent to the AgentID
type RestoreJobRequest struct {
	ID         string           `json:"id"`
	Paramiters RestoreParamiter `json:"paramiters"`
}

// RestoreRequest defines the paramiters needed for starting a Restore job
type RestoreRequest struct {
	JobID int               `json:"job"`
	From  engine.Definition `json:"from"`
	To    engine.Definition `json:"to"`
	Files []int             `json:"files"`
}

// JobFileRequest defines a file request sent to the coordinator
type JobFileRequest struct {
	File JobFile `json:"file"`
}

// JobUpdateRequest defines an update to a job
type JobUpdateRequest struct {
	State int `json:"state"`
}

// Job is a generic struct for holding backup/restore job information
type Job struct {
	ID         int         `json:"id"`
	AgentID    int         `jsong:"agentid"`
	JobType    string      `json:"type"`
	Definition interface{} `json:"definition"`
	Start      *time.Time  `json:"start"`
	End        *time.Time  `json:"end,omitempty"`
	State      int         `json:"state"`
	Message    *string     `json:"message,omitempty"`
}

// BackupDefinition Defines a backup for the Coordinator
type BackupDefinition struct {
	ID         int             `json:"id"`
	Paramiters BackupParamiter `json:"paramiters"`
	AgentID    int             `json:"agentid"`
}

// BackupPath deifines a file path to back up and any excludes to be applied
type BackupPath struct {
	Path     string   `json:"path"`
	Excludes []string `json:"excludes,omitempty"`
}

// BackupParamiter represents the paramiters for a backup operation
// Modifications to be made to the backed up files
// Engines to send the backup data to
// Paths to be backed up
type BackupParamiter struct {
	Modifications []modification.Definition `json:"modifications"`
	Engines       []engine.Definition       `json:"engines"`
	Paths         []BackupPath              `json:"paths"`
}

// RestoreParamiter represents the paramiters for a restore operation
// To represents the engine to use in the restore
// From is the engine we want to restore from
// FileSignatures is a slice of signatures to restore
// BackupDefinition is the oritinal defintion used to backup the files. This
// is needed so that the modifications can be "undone" in the reverse order and
// confirm that the "From" engine was actually used to backup the files
type RestoreParamiter struct {
	To               []engine.Definition `json:"to"`
	From             engine.Definition   `json:"from"`
	FileSignatures   []files.Signature   `json:"files"`
	BackupParamiters BackupParamiter     `json:"paramiters"`
}

// Schedule defines when a backup should run
type Schedule struct {
	ID      int    `json:"id,omitempty"`
	Backup  int    `json:"backup"`
	Seconds string `json:"seconds"`
	Minutes string `json:"minutes"`
	Hour    string `json:"hour"`
	DOM     string `json:"dom"`
	MON     string `json:"mon"`
	DOW     string `json:"dow"`
}

func (s *Schedule) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s", s.Seconds, s.Minutes, s.Hour, s.DOM, s.MON, s.DOW)
}
