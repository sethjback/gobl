package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/model"
)

// AddJobFile adds a file definition to given job id
func (d *SQLite) SaveJobFile(jobID string, f model.JobFile) error {
	split := strings.Split(f.File.Path, "/")
	jobfile, err := d.getFile(jobID, strings.Join(split[:len(split)-1], "/"), split[len(split)-1])
	if err != nil {
		if err.Error() != "Could not find that file" {
			return err
		}

		return d.insertFile(jobID, f)
	}

	return d.updatefile(jobfile.id, f)
}

func (d *SQLite) insertFile(jobID string, jf model.JobFile) error {
	values := []interface{}{jobID, jf.State}

	split := strings.Split(jf.File.Path, "/")
	values = append(values, len(split)-1)
	values = append(values, strings.Join(split[:len(split)-1], "/"))
	values = append(values, split[len(split)-1])

	var b []byte
	var e error
	if jf.Error != nil {
		b, e = json.Marshal(jf.Error)
		if e != nil {
			return errors.New("Unable to marshal error (" + e.Error() + ")")
		}
	}
	values = append(values, b)

	b, e = json.Marshal(jf.File)
	if e != nil {
		return errors.New("Unable to marshal file (" + e.Error() + ")")
	}
	values = append(values, b)

	_, err := d.Connection.Exec(
		"INSERT INTO "+filesTable+" (job, state, level, parent, name, error, file) values(?,?,?,?,?,?,?)",
		values...)

	return err
}

func (d *SQLite) updatefile(id int, jf model.JobFile) error {
	values := []interface{}{jf.State}
	split := strings.Split(jf.File.Path, "/")
	values = append(values, len(split)-1)
	values = append(values, strings.Join(split[:len(split)-1], "/"))
	values = append(values, split[len(split)-1])

	var b []byte
	var e error
	if jf.Error != nil {
		b, e = json.Marshal(jf.Error)
		if e != nil {
			return errors.New("Unable to marshal error (" + e.Error() + ")")
		}
	}
	values = append(values, b)

	b, e = json.Marshal(jf.File)
	if e != nil {
		return errors.New("Unable to marshal file (" + e.Error() + ")")
	}
	values = append(values, b)

	values = append(values, id)

	_, err := d.Connection.Exec(
		"UPDATE "+filesTable+" set state=?, level=?, parent=?, name=?, error=?, file=? where _id=?",
		values...,
	)

	return err
}

func (d *SQLite) getFile(jobID, parent, filename string) (*jobFile, error) {
	var job, state string
	var er, file []byte
	var id, level int

	err := d.Connection.QueryRow("SELECT _id, job, state, error, file, level from files WHERE job = ? and parent = ? and name = ?", jobID, parent, filename).Scan(&id, &job, &state, &er, &file, &level)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Could not find that file")
	case err != nil:
		return nil, err
	default:
		def := &jobFile{
			state: state,
			id:    id,
			level: level,
		}
		if len(er) != 0 {
			def.err = goblerr.New("", "", nil)
			err = json.Unmarshal(er, def.err)
			if err != nil {
				return nil, errors.New("Trouble unmarshalling error")
			}
		}
		if len(file) != 0 {
			err = json.Unmarshal(file, &def.file)
			if err != nil {
				return nil, errors.New("Trouble unmarshalling file")
			}
		}

		return def, nil
	}
}

func (d *SQLite) JobFiles(jobID string, filters map[string]string) ([]model.JobFile, error) {
	var wheres []string
	var vals []interface{}

	wheres = append(wheres, "job=?")
	vals = append(vals, jobID)

	for k, v := range filters {
		k = strings.ToLower(k)
		switch k {
		case "dir", "parent":
			wheres = append(wheres, "parent=?")
		case "state":
			wheres = append(wheres, "state=?")
		case "filename", "name":
			wheres = append(wheres, "name=?")
		}
		vals = append(vals, v)
	}

	ss := "SELECT error, file, state FROM " + filesTable

	if len(wheres) != 0 {
		ss += " WHERE "
		for i, w := range wheres {
			if i == 0 {
				ss += w
			} else {
				ss += " AND " + w
			}
		}
	}

	rows, err := d.Connection.Query(ss, vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobFiles []model.JobFile

	for rows.Next() {
		var err, file []byte
		var state string

		e := rows.Scan(&err, &file, &state)
		if e != nil {
			return nil, e
		}

		def := model.JobFile{State: state}

		if len(err) != 0 {
			def.Error = goblerr.New("", "", nil)
			e = json.Unmarshal(err, &def.Error)
			if e != nil {
				return nil, errors.New("unable to unmarshal error (" + e.Error() + ")")
			}
		}

		e = json.Unmarshal(file, &def.File)
		if e != nil {
			return nil, errors.New("unable to unmarshal file (" + e.Error() + ")")
		}

		jobFiles = append(jobFiles, def)
	}

	return jobFiles, nil
}

func (d *SQLite) jobFilesCount(jobID string, filters map[string]string) (int, error) {
	var wheres []string
	var vals []interface{}

	wheres = append(wheres, "job=?")
	vals = append(vals, jobID)

	for k, v := range filters {
		k = strings.ToLower(k)
		switch k {
		case "dir", "parent":
			wheres = append(wheres, "parent=?")
		case "state":
			wheres = append(wheres, "state=?")
		case "filename", "name":
			wheres = append(wheres, "name=?")
		}
		vals = append(vals, v)
	}

	ss := "SELECT COUNT(*) FROM " + filesTable

	if len(wheres) != 0 {
		ss += " WHERE "
		for i, w := range wheres {
			if i == 0 {
				ss += w
			} else {
				ss += " AND " + w
			}
		}
	}

	row := d.Connection.QueryRow(ss, vals...)
	if row == nil {
		return -1, errors.New("Query failed")
	}
	var count int
	err := row.Scan(&count)

	return count, err
}

func (d *SQLite) JobDirectories(jobID, parent string) ([]string, error) {
	split := strings.Split(parent, "/")

	rows, err := d.Connection.Query(
		"SELECT DISTINCT parent from "+filesTable+" WHERE job=? AND level=?",
		jobID, len(split)-1,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dirs := []string{}

	for rows.Next() {
		var parent string
		err = rows.Scan(&parent)
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, parent)
	}
	return dirs, nil
}
