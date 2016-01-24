package engines

import (
	"errors"
	"io"
	"strings"

	"github.com/sethjback/gobble/files"
)

// Writer provides a fan-out interface for engines
type Writer struct {
	Engines []*io.PipeWriter
}

// Option contains individual options required to configure the engines
type Option struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Default  string `json:"default"`
}

// Definition is used to identify the engine and store appropriate options
type Definition struct {
	Name    string                 `json:"name"`
	Options map[string]interface{} `json:"options"`
}

// Backup is the interface any backup engine needs to satisfy
type Backup interface {
	Backup(io.Reader, files.Signature, chan<- error)
	Retrieve(files.Signature) (io.Reader, error)
	ShouldBackup(files.Signature) (bool, error)
	Name() string
	BackupOptions() []Option
	ConfigureBackup(map[string]interface{}) error
}

// Restore is the interface any restore engine needs to satisfy
type Restore interface {
	Restore(io.Reader, files.Signature, chan<- error)
	Name() string
	RestoreOptions() []Option
	ConfigureRestore(map[string]interface{}) error
}

// A copy of the io.MultiWriter, just allows use of pipewriter
// since variadic function is eveluated at call time (i.e. even though)
// io.PipeWriter implements io.Writer, io.MultiWriter can't take a []interface{}...
func (t *Writer) Write(p []byte) (n int, err error) {
	for _, w := range t.Engines {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

// Close calls close on all the chained pipewriters
func (t *Writer) Close() {
	for i := 0; i < len(t.Engines); i++ {
		t.Engines[i].Close()
	}
}

// GetBackupEngines makes sure the agent can handle all the modifications
func GetBackupEngines(e []Definition) ([]Backup, error) {
	var es []Backup
	for _, eType := range e {
		switch strings.ToLower(eType.Name) {
		case "localfile":
			localFile := new(LocalFile)

			if err := localFile.ConfigureBackup(eType.Options); err != nil {
				return nil, err
			}

			es = append(es, localFile)

		case "logger":
			logger := new(Logger)

			if err := logger.ConfigureBackup(eType.Options); err != nil {
				return nil, err
			}

			es = append(es, logger)

		default:
			return nil, errors.New("I don't understand engine type: " + eType.Name)
		}
	}
	return es, nil
}

// GetRestoreEngine takes a definition, configures the actual engine and returns it
func GetRestoreEngine(e Definition) (Restore, error) {
	var re Restore
	switch strings.ToLower(e.Name) {
	case "localfile":
		localFile := new(LocalFile)
		if err := localFile.ConfigureRestore(e.Options); err != nil {
			return nil, err
		}
		re = localFile

	default:
		return nil, errors.New("I don't understand restore engine type: " + e.Name)

	}

	return re, nil
}
