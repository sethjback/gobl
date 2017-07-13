package engine

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/goblerr"
)

const (
	// NameLocalFile is the name of the backup engine
	NameLocalFile = "localfile"
	// LocalFileOptionSavePath is the savePath option name
	LocalFileOptionSavePath = "savePath"
	// LocalFileOptionOverwrite is the overwrite flag option name
	LocalFileOptionOverwrite = "overwrite"
	// LocalFileOptionRestorePath is the restore path option name
	LocalFileOptionRestorePath = "restorePath"

	errorAccessSavePath    = "AccessSavePathFailed"
	errorAccessRestorePath = "AccessRestorePathFailed"
)

// LocalFile backs up files to the local filesystem
type LocalFile struct {
	savePath         string
	restorePath      string
	overWrite        bool
	originalLocation bool
}

// Name returns "LocalFile"
func (e *LocalFile) Name() string {
	return NameLocalFile
}

// BackupOptions lists the available options for the bacup operation
func (e *LocalFile) SaveOptions() []Option {
	return []Option{
		Option{
			Name:     LocalFileOptionSavePath,
			Type:     "string",
			Required: true,
			Default:  ""},
		Option{
			Name:     LocalFileOptionOverwrite,
			Type:     "bool",
			Required: false,
			Default:  "true"}}
}

// ConfigureBackup options for where to save the files
func (e *LocalFile) ConfigureSave(options map[string]string) error {
	for k, v := range options {
		switch strings.ToLower(k) {

		case strings.ToLower(LocalFileOptionSavePath):
			if err := os.MkdirAll(v, 0744); err != nil {
				return goblerr.New("Configuration failed", errorAccessSavePath, fmt.Sprintf("unable to create or access %s (%s)", LocalFileOptionSavePath, err))
			}
			e.savePath = v

		case strings.ToLower(LocalFileOptionOverwrite):
			if v != "true" && v != "false" {
				return goblerr.New("Invalid option", ErrorInvalidOptionValue, LocalFileOptionOverwrite+" must true or false")
			}
			if v == "true" {
				e.overWrite = true
			}
		}
	}

	if e.savePath == "" {
		return goblerr.New("Must provide save path", ErrorRequiredOptionMissing, fmt.Sprintf("%s is required", LocalFileOptionSavePath))
	}

	return nil
}

// ShouldBackup determines if we have already saved this signature and thus wether we should save again
func (e *LocalFile) ShouldSave(file files.File) (bool, error) {
	fn, err := hashFileSig(file.Signature)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(e.savePath + string(os.PathSeparator) + fn); err == nil {
		return false, nil
	}

	return true, nil
}

// Backup saves the file to the local file system
func (e *LocalFile) Save(reader io.Reader, file files.File, errc chan<- error) {
	fn, err := hashFileSig(file.Signature)
	if err != nil {
		errc <- err
		return
	}

	if !e.overWrite {
		if _, err = os.Stat(e.savePath + string(os.PathSeparator) + fn); err == nil {
			//err is nil, which means the file exists
			errc <- errors.New("File (" + file.Path + ") exists and overWrite is false")
			return
		}
	}

	saveFile, err := os.Create(e.savePath + string(os.PathSeparator) + fn)
	if err != nil {
		errc <- err
		return
	}

	if _, err := io.Copy(saveFile, reader); err != nil {
		errc <- err
		return
	}
}

// Retrieve grabs the file from the save location and reaturns a reader to it
func (e *LocalFile) Retrieve(file files.File) (io.Reader, error) {
	fn, err := hashFileSig(file.Signature)
	if err != nil {
		return nil, err
	}

	restoreFile, err := os.OpenFile(e.savePath+string(os.PathSeparator)+fn, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	return restoreFile, nil
}

func hashFileSig(fileSig files.Signature) (string, error) {
	sig, err := json.Marshal(fileSig)
	if err != nil {
		return "", err
	}

	hash := md5.Sum(sig)
	return hex.EncodeToString(hash[:]), nil
}

// RestoreOptions lists the available options for the restore
func (e *LocalFile) RestoreOptions() []Option {
	return []Option{
		Option{
			Name:        LocalFileOptionRestorePath,
			Description: "path will be prefixed to original file path during restore. i.e. the loction where you want the files to be restored to",
			Type:        "string",
			Required:    true,
			Default:     ""},
		Option{
			Name:        LocalFileOptionOverwrite,
			Description: "whether we should overwrite the existing file if it already exists",
			Type:        "bool",
			Required:    true,
			Default:     ""}}
}

// ConfigureRestore configures the necessary options to run a local disk restore
func (e *LocalFile) ConfigureRestore(options map[string]string) error {
	oProvided := false
	rpProvided := false
	for k, v := range options {
		switch strings.ToLower(k) {
		case strings.ToLower(LocalFileOptionOverwrite):
			if v != "true" && v != "false" {
				return goblerr.New("Invalid option", ErrorInvalidOptionValue, LocalFileOptionOverwrite+" must true or false")
			}
			if v == "true" {
				e.overWrite = true
			}
			oProvided = true
		case strings.ToLower(LocalFileOptionRestorePath):
			if err := os.MkdirAll(v, 0744); err != nil {
				return goblerr.New("Configuration failed", errorAccessRestorePath, fmt.Sprintf("unable to create or access %s (%s)", LocalFileOptionRestorePath, err))
			}
			e.restorePath = v
			rpProvided = true
		}
	}

	if !rpProvided || !oProvided {
		return goblerr.New("Required options missing", ErrorRequiredOptionMissing, fmt.Sprintf("%s and %s are required", LocalFileOptionRestorePath, LocalFileOptionOverwrite))
	}

	return nil
}

// ShouldRestore checks to see if the we should restore the file
func (e *LocalFile) ShouldRestore(file files.File) (bool, error) {
	var fPath string

	if e.originalLocation {
		fPath = file.Path
	} else {
		fPath = e.restorePath + string(os.PathSeparator) + file.Path
	}

	if _, err := os.Stat(fPath); err == nil {
		return e.overWrite, nil
	} else if os.IsNotExist(err) {
		return true, nil
	} else {
		return false, err
	}
}

// Restore takes the given input stream and restores the file to the local disk
func (e *LocalFile) Restore(reader io.Reader, file files.File, errc chan<- error) {
	var fPath string
	var fFlags int

	if e.overWrite {
		fFlags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	} else {
		fFlags = os.O_CREATE | os.O_EXCL | os.O_WRONLY
	}

	if e.originalLocation {
		fPath = file.Path
	} else {

		err := os.MkdirAll(e.restorePath+string(os.PathSeparator)+file.Path, 0744)
		if err != nil {
			errc <- err
			return
		}

		fPath = e.restorePath + string(os.PathSeparator) + file.Path
	}

	rFile, err := os.OpenFile(fPath, fFlags, 0744)
	if err != nil {
		errc <- err
		return
	}
	defer rFile.Close()

	if _, err := io.Copy(rFile, reader); err != nil {
		errc <- err
		return
	}
}
