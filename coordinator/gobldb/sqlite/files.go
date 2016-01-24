package sqlite

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/sethjback/gobble/files"
	"github.com/sethjback/gobble/spec"
)

// AddJobFile adds a file definition to given job id
func (d *SQLite) AddJobFile(jobID int, f *spec.JobFile) error {
	sql := "INSERT INTO " + filesTable + " (job, state, message, path, name, signature) values(?, ?, ?, ?, ?, ?)"
	sigByte, err := json.Marshal(f.Signature)
	if err != nil {
		return err
	}
	if _, err := d.Connection.Exec(sql, jobID, f.State, f.Message, f.Signature.Path, f.Signature.Name, sigByte); err != nil {
		return err
	}
	return nil
}

// JobFiles selects a file list with the give path prefix
func (d *SQLite) JobFiles(path string, jobID int) ([]*spec.JobFile, error) {
	sql := "SELECT * FROM " + filesTable + " WHERE job=?"
	rows, err := d.Connection.Query(sql, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []*spec.JobFile

	for rows.Next() {
		var id, state int
		var message, path, name *string
		var signature, meta *[]byte
		err := rows.Scan(&id, &jobID, &state, &message, &path, &name, &signature, &meta)
		if err != nil {
			return nil, err
		}

		jf := &spec.JobFile{
			ID:    id,
			State: state}

		if message != nil {
			jf.Message = *message
		}

		err = json.Unmarshal(*signature, &jf.Signature)
		if err != nil {
			return nil, err
		}

		files = append(files, jf)
	}

	return files, nil
}

// JobFileSignatures returns a list of the file signatures in the id list
func (d *SQLite) JobFileSignatures(jobID int, fIDs []int) ([]files.Signature, error) {

	if len(fIDs) == 0 {
		return nil, errors.New("Files id list cannot be empty")
	}

	sql := "SELECT signature FROM " + filesTable + " WHERE job=? and id in(?" + strings.Repeat(",?", len(fIDs)-1) + ")"
	args := make([]interface{}, len(fIDs)+1)
	args[0] = jobID
	for i, j := range fIDs {
		args[i+1] = j
	}

	rows, err := d.Connection.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fileSigs []files.Signature

	for rows.Next() {
		var signature []byte
		err := rows.Scan(&signature)
		if err != nil {
			return nil, err
		}

		var fs files.Signature

		err = json.Unmarshal(signature, &fs)
		if err != nil {
			return nil, err
		}

		fileSigs = append(fileSigs, fs)
	}

	return fileSigs, nil
}
