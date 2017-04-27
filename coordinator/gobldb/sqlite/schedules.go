package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/sethjback/gobl/model"
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

func (d *SQLite) SaveSchedule(s model.Schedule) error {
	_, err := d.getSchedule(s.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}

		return d.insertSchedule(s)
	}

	return d.updateSchedule(s)
}

func (d *SQLite) getSchedule(id string) (*model.Schedule, error) {
	var sched []byte
	err := d.Connection.QueryRow("SELECT schedule FROM "+schedulesTable+" WHERE id=?", id).Scan(&sched)
	switch {
	case err != nil:
		return nil, err
	default:
		var def *model.Schedule
		err := json.Unmarshal(sched, &def)
		if err != nil {
			return nil, errors.New("Unable to unmarshal schedule (" + err.Error() + ")")
		}
		return def, nil
	}
}

func (d *SQLite) insertSchedule(s model.Schedule) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = d.Connection.Exec("INSERT INTO "+schedulesTable+" (id, schedule) values(?,?)", s.ID, data)
	return err
}

func (d *SQLite) updateSchedule(s model.Schedule) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = d.Connection.Exec("UPDATE "+schedulesTable+" SET schedule=? WHERE id=?", data, s.ID)
	return err
}

// ScheduleList returns a list of all saved schedules
func (d *SQLite) ScheduleList() ([]model.Schedule, error) {
	rows, err := d.Connection.Query("SELECT schedule FROM "+schedulesTable, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ss []model.Schedule

	for rows.Next() {
		var schedule []byte
		err = rows.Scan(&schedule)
		if err != nil {
			return nil, err
		}

		var s model.Schedule
		err := json.Unmarshal(schedule, &s)
		if err != nil {
			return nil, errors.New("Unable to unmarshal schedule (" + err.Error() + ")")
		}

		ss = append(ss, s)
	}

	return ss, nil
}

// DeleteSchedule removes the schedule from the DB
func (d *SQLite) DeleteSchedule(id string) error {
	_, err := d.Connection.Exec("DELETE FROM "+schedulesTable+" WHERE id=?", id)

	return err
}
