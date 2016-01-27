package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/sethjback/gobl/spec"
)

// AddBackupDefinition inserts a backup definition into the database
func (d *SQLite) AddBackupDefinition(b *spec.BackupDefinition) error {
	sql := "INSERT INTO " + definitionsTable + " (paramiters, agent) values (?, ?)"

	params, err := json.Marshal(b.Paramiters)
	if err != nil {
		return err
	}

	result, err := d.Connection.Exec(sql, params, b.AgentID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	b.ID = int(id)

	return nil
}

// UpdateBackupDefinition updates the data in the DB
func (d *SQLite) UpdateBackupDefinition(bd *spec.BackupDefinition) error {
	sql := "UPDATE " + definitionsTable + " set paramiters=? WHERE id=?"

	params, err := json.Marshal(bd.Paramiters)
	if err != nil {
		return err
	}

	if _, err = d.Connection.Exec(sql, params, bd.ID); err != nil {
		return err
	}

	return nil
}

//GetBackupDefinition gets a definition from the DB
func (d *SQLite) GetBackupDefinition(id int) (*spec.BackupDefinition, error) {

	var paramiters []byte
	var agent int
	err := d.Connection.QueryRow("SELECT * FROM "+definitionsTable+" WHERE id=?", id).Scan(&id, &agent, &paramiters)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("No backup definition with that ID")
	case err != nil:
		return nil, err
	default:
		def := &spec.BackupDefinition{ID: id, AgentID: agent}
		if err := json.Unmarshal(paramiters, &def.Paramiters); err != nil {
			return nil, err
		}
		return def, nil
	}

}

// BackupDefinitionList pulls a list of all BackupDefintions from the DB
func (d *SQLite) BackupDefinitionList() ([]*spec.BackupDefinition, error) {
	rows, err := d.Connection.Query("SELECT * from "+definitionsTable, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Done this way so that the api returns empty json array rather than nil
	var b []*spec.BackupDefinition
	for rows.Next() {
		var paramiters []byte
		var id, agent int
		err := rows.Scan(&id, &agent, &paramiters)
		if err != nil {
			return nil, err
		}

		def := &spec.BackupDefinition{ID: id, AgentID: agent}
		if err := json.Unmarshal(paramiters, &def.Paramiters); err != nil {
			return nil, err
		}

		b = append(b, def)
	}

	return b, nil

}

// DeleteBackupDefinition removes a definition from the DB
func (d *SQLite) DeleteBackupDefinition(id int) error {
	if _, err := d.Connection.Exec("DELETE FROM "+definitionsTable+" WHERE id=?", id); err != nil {
		return err
	}

	return nil
}
