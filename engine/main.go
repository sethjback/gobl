package engine

import (
	"errors"
	"io"
	"strings"

	"github.com/sethjback/gobl/files"
)

const (
	// ErrorInvalidOptionValue is used by engines to indicate a provided option value was invalid
	ErrorInvalidOptionValue = "InvalidOptionValue"
	// ErrorRequiredOptionMissing is used by engines to indicate a required option is absent
	ErrorRequiredOptionMissing = "RequiredOptionMissing"
)

// Option contains individual options required to configure the engines
type Option struct {
	// Name of the option
	Name string `json:"name"`
	// Description of what the option does
	Description string `json:"description"`
	// Type of value for this option
	Type string `json:"type"`
	// Required indicates if this option is non-optional
	Required bool `json:"required"`
	// Default value for this option
	Default interface{} `json:"default"`
}

// Definition is used to identify the engine and store appropriate options
// The definitions are used during jobs to store engine configuration details
type Definition struct {
	// Name of the engine
	Name string `json:"name"`
	// Options used to configure it
	Options map[string]string `json:"options"`
}

// Saver is the interface an engine needs to satisfy for saving (backing up) data
type Saver interface {
	// Save process the input bytes. Signauture gives information about the file being saved
	// any errors encountered during the write should be sent over errc
	Save(input io.Reader, signature files.File, errc chan<- error)
	// Retrieve the file represented by signature
	Retrieve(signature files.File) (io.Reader, error)
	// ShouldBackup indicates whether the engine needs to process the file
	ShouldSave(signature files.File) (bool, error)
	// Name of the backup engine
	Name() string
	// BackupOptions available to configure the engine
	SaveOptions() []Option
	// ConfigureBackup sets the appropriate options
	ConfigureSave(options map[string]string) error
}

// Restorer is the interface an engine needs to satisfy for restoring data
type Restorer interface {
	// Restore processes the input. Signature gives information aobut the file
	// Any errors encountered when restoring the file should be sent over errc
	Restore(input io.Reader, signature files.File, errc chan<- error)
	// ShouldRestore indicates whether the engine needs to process the file
	ShouldRestore(files.File) (bool, error)
	// Name of the restore engine
	Name() string
	// RestoreOptions available to configure the engine
	RestoreOptions() []Option
	// ConfigureRestore sets the appropriate options
	ConfigureRestore(map[string]string) error
}

// BuildSavers returns a slice of configured savers
func BuildSavers(definitions []Definition) ([]Saver, error) {
	var sers []Saver
	for _, d := range definitions {
		switch strings.ToLower(d.Name) {
		case NameLocalFile:
			localFile := &LocalFile{}

			if err := localFile.ConfigureSave(d.Options); err != nil {
				return nil, err
			}

			sers = append(sers, localFile)

		case NameLogger:
			logger := &Logger{}

			if err := logger.ConfigureSave(d.Options); err != nil {
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

		case NameLogger:
			logger := &Logger{}
			if err := logger.ConfigureRestore(d.Options); err != nil {
				return nil, err
			}
			rers = append(rers, logger)

		default:
			return nil, errors.New("I don't understand restore engine type: " + d.Name)

		}
	}

	return rers, nil
}
