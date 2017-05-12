package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/sethjback/gobl/gobldb"
	"github.com/sethjback/gobl/model"
)

func (d *SQLite) SaveJob(job model.Job) error {
	_, err := d.GetJob(job.ID)
	if err != nil {
		if err.Error() != "No job with that ID" {
			return err
		}

		return d.insertJob(job)
	}

	return d.updateJob(job)
}

func (d *SQLite) updateJob(job model.Job) error {
	job.Agent = nil

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = d.Connection.Exec(
		"UPDATE "+jobsTable+" SET start=?, end=?, state=?, data=? WHERE id=?",
		job.Meta.Start.Unix(), job.Meta.End.Unix(), job.Meta.State, data, job.ID)

	return err
}

func (d *SQLite) insertJob(job model.Job) error {
	sID, err := d.getAgentSQLId(job.Agent.ID)
	if err != nil {
		return err
	}

	job.Agent = nil

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = d.Connection.Exec(
		"INSERT INTO "+jobsTable+" (id, agent, start, end, state, data) values(?,?,?,?,?,?)",
		job.ID, sID, job.Meta.Start.Unix(), job.Meta.End.Unix(), job.Meta.State, data)

	return err
}

func (d *SQLite) GetJob(id string) (*model.Job, error) {
	var data []byte
	var agentID int
	err := d.Connection.QueryRow("SELECT agent, data FROM "+jobsTable+" WHERE id=?", id).Scan(&agentID, &data)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("No job with that ID")
	case err != nil:
		return nil, err
	default:
		var def *model.Job
		if err = json.Unmarshal(data, &def); err != nil {
			return nil, err
		}

		if def.Agent, err = d.getAgent(agentID); err != nil {
			return nil, err
		}

		def.Meta.Errors, _ = d.jobFilesCount(id, map[string]string{"state": model.StateFailed})
		def.Meta.Complete, _ = d.jobFilesCount(id, map[string]string{"state": model.StateFinished})
		def.Meta.Total, _ = d.jobFilesCount(id, nil)

		return def, err
	}
}

func (d *SQLite) GetJobs(filters map[string]string) ([]model.Job, error) {
	var wheres []string
	var vals []interface{}
	limit := 10
	offset := 0

	for k, v := range filters {
		k = strings.ToLower(k)
		switch k {
		case "state":
			wheres = append(wheres, "state=?")
			vals = append(vals, v)
		case "start":
			i, err := parseDate(v)
			if err != nil {
				return nil, errors.New("start: " + err.Error())
			}
			wheres = append(wheres, "start > ?")
			vals = append(vals, i)
		case "end":
			i, err := parseDate(v)
			if err != nil {
				return nil, errors.New("end: " + err.Error())
			}
			wheres = append(wheres, "end < ?")
			vals = append(vals, i)
		case "agentid":
			aid, err := d.getAgentSQLId(v)
			if err != nil {
				return nil, errors.New("Invalid agentId: " + err.Error())
			}
			wheres = append(wheres, "agent=?")
			vals = append(vals, aid)
		case "limit":
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("invalid limit")
			}
			limit = i
		case "offset":
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("invalid offset")
			}
			offset = i
		}
	}

	sql := "SELECT id FROM " + jobsTable
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
	var jobIDs []string

	for rows.Next() {
		var id string

		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		jobIDs = append(jobIDs, id)
	}
	rows.Close()

	jobs := make([]model.Job, len(jobIDs))

	for i, ID := range jobIDs {
		j, err := d.GetJob(ID)
		if err != nil {
			return nil, err
		}
		jobs[i] = *j
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
