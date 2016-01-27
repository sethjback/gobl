package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/sethjback/gobl/spec"
)

// AddSchedule saves a schedule to the DB
func (d *SQLite) AddSchedule(s *spec.Schedule) error {

	sql := "INSERT INTO " + schedulesTable + " (backup, schedule) values(?, ?)"

	sched, err := json.Marshal(s)
	if err != nil {
		return errors.New("Could not marshal schedule")
	}

	result, err := d.Connection.Exec(sql, s.Backup, sched)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil
	}

	s.ID = int(id)

	return nil
}

// ScheduleList returns a list of all saved schedules
func (d *SQLite) ScheduleList() ([]*spec.Schedule, error) {
	rows, err := d.Connection.Query("SELECT * FROM "+schedulesTable, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ss []*spec.Schedule

	for rows.Next() {
		var id, backup int
		var body []byte
		err = rows.Scan(&id, &backup, &body)
		if err != nil {
			return nil, err
		}

		var s spec.Schedule
		json.Unmarshal(body, &s)

		s.Backup = backup
		s.ID = id
		ss = append(ss, &s)
	}

	return ss, nil
}

// GetSchedule returns single schedule
func (d *SQLite) GetSchedule(id int) (*spec.Schedule, error) {
	var s []byte
	var backup int
	err := d.Connection.QueryRow("SELECT * FROM "+schedulesTable+" WHERE id=?", id).Scan(&id, &backup, &s)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("No schedule with that ID")
	case err != nil:
		return nil, err
	default:
		var sched spec.Schedule
		json.Unmarshal(s, &sched)
		sched.ID = id
		return &sched, nil
	}
}

// DeleteSchedule removes the schedule from the DB
func (d *SQLite) DeleteSchedule(id int) error {
	if _, err := d.Connection.Exec("DELETE FROM "+schedulesTable+" WHERE id=?", id); err != nil {
		return err
	}

	return nil
}

// UpdateSchedule completely replaces the schedule
func (d *SQLite) UpdateSchedule(s *spec.Schedule) error {
	sql := "UPDATE " + schedulesTable + " set backup=?, schedule=? WHERE id=?"

	sbyte, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if _, err = d.Connection.Exec(sql, s.Backup, sbyte, s.ID); err != nil {
		return err
	}
	return nil
}
