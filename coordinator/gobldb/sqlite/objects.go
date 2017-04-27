package sqlite

import (
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/goblerr"
)

type job struct {
	id      int
	agentID int
}

type jobFile struct {
	id     int
	job    string
	state  string
	err    goblerr.Error
	file   files.File
	level  int
	parent string
	name   string
}
