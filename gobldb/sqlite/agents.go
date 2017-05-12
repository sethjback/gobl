package sqlite

import (
	"database/sql"
	"errors"

	"github.com/sethjback/gobl/model"
)

func (d *SQLite) SaveAgent(a model.Agent) error {
	_, err := d.GetAgent(a.ID)
	if err != nil {
		if err.Error() != "No agent with that ID" {
			return err
		}

		return insertAgent(d, a)
	}

	return updateAgent(d, a)
}

func updateAgent(d *SQLite, a model.Agent) error {
	sql := "UPDATE " + agentsTable + " set name=?, address=?, publicKey=? WHERE id=?"
	_, err := d.Connection.Exec(sql, a.Name, a.Address, a.PublicKey, a.ID)
	return err
}

func insertAgent(d *SQLite, a model.Agent) error {
	sql := "INSERT INTO " + agentsTable + " (id, name, address, publicKey) values (?, ?, ?, ?)"

	_, err := d.Connection.Exec(sql, a.ID, a.Name, a.Address, a.PublicKey)

	return err
}

func (d *SQLite) getAgent(id int) (*model.Agent, error) {
	var sID, name, address, publicKey string
	err := d.Connection.QueryRow("SELECT id, name, address, publicKey FROM "+agentsTable+" WHERE _id=?", id).Scan(&sID, &name, &address, &publicKey)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("No agent with that ID")
	case err != nil:
		return nil, err
	default:
		return &model.Agent{ID: sID, Name: name, Address: address, PublicKey: publicKey}, nil
	}
}

func (d *SQLite) getAgentSQLId(id string) (int, error) {
	var sID int
	err := d.Connection.QueryRow("SELECT _id FROM "+agentsTable+" WHERE id=?", id).Scan(&sID)
	return sID, err
}

// AgentList returns a list of all agents in the DB
func (d *SQLite) AgentList() ([]model.Agent, error) {
	rows, err := d.Connection.Query("SELECT id, name, address, publicKey from "+agentsTable, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var agents []model.Agent

	for rows.Next() {
		var ID, name, address, publicKey string
		err := rows.Scan(&ID, &name, &address, &publicKey)
		if err != nil {
			return nil, err
		}
		agents = append(agents, model.Agent{ID: ID, Name: name, Address: address, PublicKey: publicKey})
	}

	return agents, nil
}

// GetAgent retrieves an agent from the DB
func (d *SQLite) GetAgent(id string) (*model.Agent, error) {
	var name, address, publicKey string
	err := d.Connection.QueryRow("SELECT name, address, publicKey FROM "+agentsTable+" WHERE id=?", id).Scan(&name, &address, &publicKey)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("No agent with that ID")
	case err != nil:
		return nil, err
	default:
		return &model.Agent{ID: id, Name: name, Address: address, PublicKey: publicKey}, nil
	}
}
