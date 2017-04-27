package model

import (
	"time"

	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/modification"
)

type Job struct {
	ID         string         `json:"id"`
	Agent      *Agent         `json:"agent,omitempty"`
	Definition *JobDefinition `json:"definition"`
	Meta       *JobMeta       `json:"meta"`
}

type JobDefinition struct {
	ID            string                    `json:"id,omitempty"`
	To            []engine.Definition       `json:"to"`
	From          *engine.Definition        `json:"from,omitempty"`
	Modifications []modification.Definition `json:"modifications"`
	Paths         []Path                    `json:"paths,omitempty"`
	Files         []files.File              `json:"files,omitempty"`
}

type JobMeta struct {
	State    string    `json:"state"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration float64   `json:"duration"`
	Total    int       `json:"total"`
	Complete int       `json:"complete"`
	Skipped  int       `json:"skipped"`
	Errors   int       `json:"errors"`
}

type JobFile struct {
	File  files.File    `json:"file"`
	State string        `json:"state"`
	Error goblerr.Error `json:"error,omitempty"`
}

type Path struct {
	Root     string   `json:"root"`
	Excludes []string `json:"excludes"`
}
