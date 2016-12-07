package engines

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/sethjback/gobl/files"
)

const NameLogger = "logger"

// Logger is an engine that takes incoming backup requests and simply logs them
// to a file
type Logger struct {
	logPath   string
	overWrite bool
}

type logLine struct {
	Start         string   `json:"start"`
	End           string   `json:"end"`
	File          string   `json:"file"`
	Modifications []string `json:"modifications"`
	Engines       []string `json:"engines"`
	Signature     string   `json:"signature"`
	Chunks        int64    `json:"chunks"`
	Bytes         int64    `json:"bytes"`
}

// Name return "Logger"
func (e *Logger) Name() string {
	return "Logger"
}

// BackupOptions lists all the availalbe options
func (e *Logger) BackupOptions() []Option {
	return []Option{
		Option{
			Name:     "logPath",
			Type:     "string",
			Required: true,
			Default:  ""},
		Option{
			Name:     "overWrite",
			Type:     "bool",
			Required: false,
			Default:  "false"}}
}

// ConfigureBackup configures the engine so it is ready to use
func (e *Logger) ConfigureBackup(options map[string]interface{}) error {
	val, ok := options["logPath"]
	if !ok {
		return errors.New("logger: Invalid options: need logPath")
	}

	vstring, ok := val.(string)
	if !ok {
		return errors.New("logger: Invalid logPath: must be string")
	}

	e.logPath = vstring

	val, ok = options["overWrite"]
	if !ok {
		e.overWrite = false
	} else {
		vbool, ok := val.(bool)
		if !ok {
			return errors.New("logger: Invalid overWrite: must be bool")
		}

		e.overWrite = vbool
	}

	var lFlags int

	if e.overWrite {
		lFlags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	} else {
		lFlags = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	}
	lFile, err := os.OpenFile(e.logPath, lFlags, 0600)
	if err != nil {
		return errors.New("logger: Could not open Logfile for writing")
	}
	lFile.Close()

	return nil

}

// ShouldBackup always returns true: we want to log each file
func (e *Logger) ShouldBackup(fileSig files.Signature) (bool, error) {
	return true, nil
}

// Backup collects information about the file then writes it to a log once the backup has finished
func (e *Logger) Save(reader io.Reader, fileSig files.Signature, errc chan<- error) {
	fileData := new(logLine)
	fileData.File = fileSig.Path + "/" + fileSig.Name
	fileData.Modifications = fileSig.Modifications
	sig, _ := json.Marshal(fileSig)
	fileData.Signature = string(sig)
	fileData.Start = time.Now().String()

	buf := make([]byte, 0, 4*1024)
	for {
		n, err := reader.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			errc <- err
			return
		}
		fileData.Chunks++
		fileData.Bytes += int64(len(buf))
	}

	fileData.End = time.Now().String()

	dataBytes, _ := json.Marshal(fileData)

	lfile, err := os.OpenFile(e.logPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		errc <- err
	}
	defer lfile.Close()

	lfile.WriteString(string(dataBytes) + "\n")
}

// Retrieve does nothing
func (e *Logger) Retrieve(fileSig files.Signature) (io.Reader, error) {
	return nil, errors.New("Cannot use Logger to restore files")
}
