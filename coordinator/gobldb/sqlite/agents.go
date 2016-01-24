package sqlite

import (
	"database/sql"
	"errors"

	"github.com/sethjback/gobble/spec"
)

// AddAgent inserts an agent into the DB
func (d *SQLite) AddAgent(s *spec.Agent) error {
	sql := "INSERT INTO " + agentsTable + " (name, address, publicKey) values (?, ?, ?)"

	result, err := d.Connection.Exec(sql, s.Name, s.Address, s.PublicKey)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = int(id)
	return nil
}

// AgentList returns a list of all agents in the DB
func (d *SQLite) AgentList() ([]*spec.Agent, error) {
	rows, err := d.Connection.Query("SELECT * from "+agentsTable, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var agents []*spec.Agent

	for rows.Next() {
		var id int
		var name, address, publicKey string
		err := rows.Scan(&id, &name, &address, &publicKey)
		if err != nil {
			return nil, err
		}
		agents = append(agents, &spec.Agent{id, name, address, publicKey})
	}
	return agents, nil
}

// GetAgent retrieves an agent from the DB
func (d *SQLite) GetAgent(id int) (*spec.Agent, error) {
	var name, address, publicKey string
	err := d.Connection.QueryRow("SELECT * FROM "+agentsTable+" WHERE id=?", id).Scan(&id, &name, &address, &publicKey)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("No agent with that ID")
	case err != nil:
		return nil, err
	default:
		return &spec.Agent{id, name, address, publicKey}, nil
	}
}

// UpdateAgent puts the provided spec into the DB, overriding all old data
func (d *SQLite) UpdateAgent(agent *spec.Agent) error {
	sql := "UPDATE " + agentsTable + " set name=?, address=?, publicKey=? WHERE id=?"
	_, err := d.Connection.Exec(sql, agent.Name, agent.Address, agent.PublicKey, agent.ID)
	if err != nil {
		return err
	}

	return nil
}
