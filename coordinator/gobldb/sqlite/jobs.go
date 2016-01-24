package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/sethjback/gobble/coordinator/gobldb"
	"github.com/sethjback/gobble/spec"
)

// CreateBackupJob inserts a new backup job record into the DB and marks it Running
func (d *SQLite) CreateBackupJob(b *spec.BackupDefinition) (*spec.Job, error) {
	sql := "INSERT INTO " + jobsTable + " (jobtype, agent, definition, start, state) values(?, ?, ?, ?, ?)"

	t := time.Now()
	bs, err := json.Marshal(b.Paramiters)
	if err != nil {
		return nil, err
	}

	result, err := d.Connection.Exec(sql, "backup", b.AgentID, bs, t, spec.Running)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &spec.Job{ID: int(id), JobType: "backup", Definition: b, Start: &t, State: spec.Running}, nil
}

// CreateRestoreJob inserts a restore job into the DB
func (d *SQLite) CreateRestoreJob(agentID int, r *spec.RestoreRequest) (*spec.Job, error) {
	sql := "INSERT INTO " + jobsTable + " (jobtype, agent, definition, start, state) values(?, ?, ?, ?, ?)"

	t := time.Now()
	rr, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	result, err := d.Connection.Exec(sql, "restore", agentID, rr, t, spec.Running)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &spec.Job{ID: int(id), JobType: "restore", Definition: r, Start: &t, State: spec.Running}, nil
}

// GetJob returns the given job spec
func (d *SQLite) GetJob(id int) (*spec.Job, error) {
	var start, end *time.Time
	var state, agent int
	var definition []byte
	var message, jobType *string
	err := d.Connection.QueryRow("SELECT * from "+jobsTable+" WHERE id=?", id).Scan(&id, &jobType, &agent, &definition, &start, &end, &state, &message)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("No backup job with that ID")
	case err != nil:
		return nil, err
	default:
		jobSpec := &spec.Job{
			ID:      id,
			AgentID: agent,
			JobType: *jobType,
			Start:   start,
			End:     end,
			State:   state,
			Message: message}
		if *jobType == "backup" {
			var bp spec.BackupParamiter
			json.Unmarshal(definition, &bp)
			jobSpec.Definition = bp
		} else {
			var rr spec.RestoreRequest
			json.Unmarshal(definition, &rr)
			jobSpec.Definition = rr
		}
		return jobSpec, nil
	}
}

// UpdateJob updates a job's state, endtime, and message
func (d *SQLite) UpdateJob(j *spec.Job) error {
	sql := "UPDATE " + jobsTable + " set state=?, end=?, message=? WHERE id=?"
	_, err := d.Connection.Exec(sql, j.State, j.End, j.Message, j.ID)
	if err != nil {
		return err
	}

	return nil
}

// JobErrorCount returns the number of errors for a particular job
func (d *SQLite) JobErrorCount(jobID int) (int, error) {
	sql := "SELECT COUNT(id) FROM " + filesTable + " WHERE state=? AND job=?"
	var count int
	err := d.Connection.QueryRow(sql, spec.Errors, jobID).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// JobQuery runs a query against the jobs table and returns a list of jobs conforming to the paramiters
func (d *SQLite) JobQuery(params map[string]string) ([]*spec.Job, error) {

	var wheres []string
	var vals []interface{}
	limit := 10
	offset := 0

	for key, value := range params {
		switch key {
		case "state":
			i, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("Invalid state flag")
			}
			wheres = append(wheres, "state=?")
			vals = append(vals, i)
		case "start":
			i, err := parseDate(value)
			if err != nil {
				return nil, errors.New("start: " + err.Error())
			}
			wheres = append(wheres, "start > datetime(?, 'unixepoch', 'localtime')")
			vals = append(vals, i)
		case "end":
			i, err := parseDate(value)
			if err != nil {
				return nil, errors.New("end: " + err.Error())
			}
			wheres = append(wheres, "end < datetime(?, 'unixepoch', 'localtime')")
			vals = append(vals, i)
		case "agentid":
			i, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("Invalid agentid")
			}
			wheres = append(wheres, "agent=?")
			vals = append(vals, i)
		case "limit":
			i, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("invalid limit")
			}
			limit = i
		case "offset":
			i, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("invalid offset")
			}
			offset = i
		}
	}
	sql := "SELECT * FROM " + jobsTable
	if len(wheres) != 0 {
		sql += " WHERE "
		for i, w := range wheres {
			if i == 0 {
				sql += w
			} else {
				sql += " AND " + w
			}
		}
	}

	sql += " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)

	rows, err := d.Connection.Query(sql, vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []*spec.Job
	for rows.Next() {
		var jobType, message *string
		var state, id, agent int
		var definition []byte
		var start, end *time.Time

		err := rows.Scan(&id, &jobType, &agent, &definition, &start, &end, &state, &message)
		if err != nil {
			return nil, err
		}

		jobSpec := &spec.Job{
			ID:      id,
			AgentID: agent,
			JobType: *jobType,
			Start:   start,
			End:     end,
			State:   state,
			Message: message}
		if *jobType == "backup" {
			var bp spec.BackupParamiter
			json.Unmarshal(definition, &bp)
			jobSpec.Definition = bp
		} else if *jobType == "restore" {
			var rr spec.RestoreRequest
			json.Unmarshal(definition, &rr)
			jobSpec.Definition = rr
		}
		jobs = append(jobs, jobSpec)
	}

	return jobs, nil
}

func parseDate(date string) (int, error) {

	//attempt to parseDate
	t, err := time.Parse(gobldb.QueryDateFormat, date)
	if err == nil {
		return int(t.Unix()), nil
	}

	di, err := strconv.Atoi(date)
	if err == nil {
		return di, nil
	}

	return -1, errors.New("Invalid time provided. Must be unix timestamp or in format yyyy-mm-dd hh:mm")
}
