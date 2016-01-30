package sqlite

import (
	"database/sql"

	//use the sqlite3 interface for go
	_ "github.com/mattn/go-sqlite3"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/util/log"
)

// SQLite DB Implementation
type SQLite struct {
	Connection *sql.DB
}

// Init opens the DB and creates tables if necessary
func (d *SQLite) Init(options config.DB) error {

	if len(options.DBPath) == 0 {
		log.Warnf("sqlite", "DB Path empty, this will create in-memory db: probably not what you wanted!")
	}

	db, err := sql.Open("sqlite3", options.DBPath)
	if err != nil {
		return err
	}
	d.Connection = db

	_, err = d.Connection.Exec(createAgentsTable, "")
	if err != nil {
		return err
	}
	_, err = d.Connection.Exec(createDefinitionsTable, "")
	if err != nil {
		return err
	}
	_, err = d.Connection.Exec(createFilesTable, "")
	if err != nil {
		return err
	}
	_, err = d.Connection.Exec(createJobsTable, "")
	if err != nil {
		return err
	}
	_, err = d.Connection.Exec(createSchedulesTable, "")
	if err != nil {
		return err
	}

	return nil
}
