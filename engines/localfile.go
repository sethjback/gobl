package engines

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/sethjback/gobble/files"
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
	return "LocalFile"
}

// BackupOptions lists the available options for the bacup operation
func (e *LocalFile) BackupOptions() []Option {
	return []Option{
		Option{
			Name:     "savePath",
			Type:     "string",
			Required: true,
			Default:  ""},
		Option{
			Name:     "overWrite",
			Type:     "bool",
			Required: false,
			Default:  "true"}}
}

// ConfigureBackup options for where to save the files
func (e *LocalFile) ConfigureBackup(options map[string]interface{}) error {
	val, ok := options["savePath"]
	if !ok {
		return errors.New("Invalid options: need savePath")
	}

	vstring, ok := val.(string)
	if !ok {
		return errors.New("Invalid savePath: must be string")
	}

	val, ok = options["overWrite"]
	if !ok {
		return errors.New("Invalid Options: need overWrite flag")
	}

	vbool, ok := val.(bool)
	if !ok {
		return errors.New("Invalid overWrite: must be bool")
	}

	e.overWrite = vbool

	if err := os.MkdirAll(vstring, 0744); err != nil {
		return errors.New("Could not create or access savePath")
	}

	e.savePath = vstring

	return nil
}

// ShouldBackup determines if we have already saved this signature and thus wether we should save again
func (e *LocalFile) ShouldBackup(fileSig files.Signature) (bool, error) {
	sig, err := json.Marshal(fileSig)
	if err != nil {
		return false, err
	}

	//File name is the md5Hash of the file signature
	hash := md5.Sum([]byte(sig))
	fn := hex.EncodeToString(hash[:])
	if _, err := os.Stat(e.savePath + string(os.PathSeparator) + fn); err == nil {
		return false, nil
	}

	return true, nil
}

// Backup saves the file to the local file system
func (e *LocalFile) Backup(reader io.Reader, fileSig files.Signature, errc chan<- error) {
	sig, err := json.Marshal(fileSig)
	if err != nil {
		errc <- err
		return
	}

	//File name is the md5Hash of the file signature
	hash := md5.Sum([]byte(sig))
	fn := hex.EncodeToString(hash[:])

	if !e.overWrite {
		if _, err := os.Stat(e.savePath + string(os.PathSeparator) + fn); err == nil {
			//err is nil, which means the file exists
			errc <- errors.New("File (" + fileSig.Path + string(os.PathSeparator) + fileSig.Name + ") exists and overWrite is false")
			return
		}
	}

	file, err := os.Create(e.savePath + string(os.PathSeparator) + fn)
	if err != nil {
		errc <- err
		return
	}

	if _, err := io.Copy(file, reader); err != nil {
		errc <- err
		return
	}
}

// Retrieve grabs the file from the save location and reaturns a reader to it
func (e *LocalFile) Retrieve(fileSig files.Signature) (io.Reader, error) {
	sig, err := json.Marshal(fileSig)
	if err != nil {
		return nil, err
	}
	hash := md5.Sum([]byte(sig))
	fn := hex.EncodeToString(hash[:])

	file, err := os.OpenFile(e.savePath+string(os.PathSeparator)+fn, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// RestoreOptions lists the available options for the bacup operation
func (e *LocalFile) RestoreOptions() []Option {
	return []Option{
		Option{
			Name:     "restorePath",
			Type:     "string",
			Required: true,
			Default:  ""},
		Option{
			Name:     "overWrite",
			Type:     "bool",
			Required: true,
			Default:  ""}}
}

// ConfigureRestore configures the necessary options to run a local disk restore
func (e *LocalFile) ConfigureRestore(options map[string]interface{}) error {
	val, ok := options["overWrite"]
	if !ok {
		return errors.New("Invalid Options: need overWrite flag")
	}

	vbool, ok := val.(bool)
	if !ok {
		return errors.New("Invalid overWrite: must be bool")
	}

	e.overWrite = vbool

	val, ok = options["restorePath"]
	if !ok {
		return errors.New("Invalid options: need restorePath")
	}

	vstring, ok := val.(string)
	if !ok {
		return errors.New("Invalid restorePath: must be string")
	}

	e.restorePath = vstring

	if vstring == "original" {
		e.originalLocation = true
	} else {
		e.originalLocation = false

		if err := os.MkdirAll(vstring, 0744); err != nil {
			return errors.New("Could not create or access restorePath")
		}
	}

	return nil
}

// Restore takes the given input stream and restores the file to the local disk
func (e *LocalFile) Restore(reader io.Reader, fileSig files.Signature, errc chan<- error) {
	var fPath string
	var fFlags int

	if e.overWrite {
		fFlags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	} else {
		fFlags = os.O_CREATE | os.O_EXCL | os.O_WRONLY
	}

	if e.originalLocation {
		fPath = fileSig.Path + string(os.PathSeparator) + fileSig.Name
	} else {

		err := os.MkdirAll(e.restorePath+string(os.PathSeparator)+fileSig.Path, 0744)
		if err != nil {
			errc <- err
			return
		}

		fPath = e.restorePath + string(os.PathSeparator) + fileSig.Path + string(os.PathSeparator) + fileSig.Name
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
