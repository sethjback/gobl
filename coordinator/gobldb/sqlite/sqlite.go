package sqlite

import (
	"database/sql"
	"errors"

	//use the sqlite3 interface for go
	_ "github.com/mattn/go-sqlite3"
)

// SQLite DB Implementation
type SQLite struct {
	Connection *sql.DB
}

// Init opens the DB and creates tables if necessary
func (d *SQLite) Init(options map[string]interface{}) error {
	//sql.Register(driverName, &sqlite3.SQLiteDriver{})
	path, ok := options["DBPath"].(string)
	if !ok {
		return errors.New("Sqlite requires db path")
	}

	db, err := sql.Open("sqlite3", path)
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
