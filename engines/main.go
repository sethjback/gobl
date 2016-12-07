package engines

import (
	"errors"
	"io"
	"strings"

	"github.com/sethjback/gobl/files"
)

const (
	OperationBackup  = 1
	OperationRestore = 2
)

// Option contains individual options required to configure the engines
type Option struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default"`
}

// Definition is used to identify the engine and store appropriate options
type Definition struct {
	Name    string                 `json:"name"`
	Options map[string]interface{} `json:"options"`
}

// Backup is the interface any backup engine needs to satisfy
type Saver interface {
	// Save process the input bytes. Signauture gives information about the file being saved
	// any errors encountered during the write should be sent over errc
	Save(input io.Reader, signature files.Signature, errc chan<- error)
	// Retrieve the file represented by signature
	Retrieve(signature files.Signature) (io.Reader, error)
	// ShouldBackup indicates whether the engine needs to process the file
	ShouldBackup(files.Signature) (bool, error)
	// Name of the backup engine
	Name() string
	// BackupOptions available to configure the engine
	BackupOptions() []Option
	// ConfigureBackup sets the appropriate options
	ConfigureBackup(map[string]interface{}) error
}

// Restore is the interface any restore engine needs to satisfy
type Restorer interface {
	// Restore processes the input. Signature gives information aobut the file
	// Any errors encountered when restoring the file should be sent over errc
	Restore(input io.Reader, signature files.Signature, errc chan<- error)
	// ShouldRestore indicates whether the engine needs to process the file
	ShouldRestore(files.Signature) (bool, error)
	// Name of the restore engine
	Name() string
	// RestoreOptions available to configure the engine
	RestoreOptions() []Option
	// ConfigureRestore sets the appropriate options
	ConfigureRestore(map[string]interface{}) error
}

// BuildSavers returns a slice of configured savers
func BuildSavers(definitions []Definition) ([]Saver, error) {
	var sers []Saver
	for _, d := range definitions {
		switch strings.ToLower(d.Name) {
		case NameLocalFile:
			localFile := &LocalFile{}

			if err := localFile.ConfigureBackup(d.Options); err != nil {
				return nil, err
			}

			sers = append(sers, localFile)

		case NameLogger:
			logger := &Logger{}

			if err := logger.ConfigureBackup(d.Options); err != nil {
				return nil, err
			}

			sers = append(sers, logger)

		default:
			return nil, errors.New("I don't understand engine type: " + d.Name)
		}
	}
	return sers, nil
}

// BuildRestorers returns a slice of configured restorers
func BuildRestorers(definitions []Definition) ([]Restorer, error) {
	var rers []Restorer
	for _, d := range definitions {
		switch strings.ToLower(d.Name) {
		case NameLocalFile:
			localFile := &LocalFile{}
			if err := localFile.ConfigureRestore(d.Options); err != nil {
				return nil, err
			}
			rers = append(rers, localFile)

		default:
			return nil, errors.New("I don't understand restore engine type: " + d.Name)

		}
	}

	return rers, nil
}
