package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/sethjback/gobl/model"
)

func (d *SQLite) SaveJobDefinition(jd model.JobDefinition) error {
	_, err := d.GetJobDefinition(jd.ID)
	if err != nil {
		if err.Error() != "No job definition with that ID" {
			return err
		}
		return insertJobDef(d, jd)
	}

	return updateJobDef(d, jd)
}

func insertJobDef(d *SQLite, jd model.JobDefinition) error {
	sql := "INSERT INTO " + definitionsTable + " (id, data) values (?, ?)"

	data, err := json.Marshal(jd)
	if err != nil {
		return err
	}

	_, err = d.Connection.Exec(sql, jd.ID, data)
	return err
}

func updateJobDef(d *SQLite, jd model.JobDefinition) error {
	data, err := json.Marshal(jd)
	if err != nil {
		return err
	}

	sql := "UPDATE " + definitionsTable + " set data=? WHERE id=?"
	_, err = d.Connection.Exec(sql, data, jd.ID)
	return err
}

func (d *SQLite) GetJobDefinition(id string) (*model.JobDefinition, error) {
	var data []byte

	err := d.Connection.QueryRow("SELECT data FROM "+definitionsTable+" WHERE id=?", id).Scan(&data)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("No job definition with that ID")
	case err != nil:
		return nil, err
	default:
		var def *model.JobDefinition
		if err = json.Unmarshal(data, &def); err != nil {
			return nil, err
		}

		return def, err
	}
}

func (d *SQLite) DeleteJobDefinition(id string) error {
	_, err := d.Connection.Exec("DELETE FROM "+definitionsTable+" WHERE id=?", id)

	return err
}

func (d *SQLite) GetJobDefinitions() ([]model.JobDefinition, error) {
	rows, err := d.Connection.Query("SELECT id, data from "+definitionsTable, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var defs []model.JobDefinition
	for rows.Next() {
		var def model.JobDefinition
		var data []byte
		var id string
		if err = rows.Scan(&id, &data); err != nil {
			return nil, err
		}

		if err = json.Unmarshal(data, &def); err != nil {
			def = model.JobDefinition{ID: id}
		}

		defs = append(defs, def)
	}

	return defs, nil
}
