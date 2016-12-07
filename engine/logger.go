package engine

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/goblerr"
)

const (
	NameLogger            = "logger"
	LoggerOptionLogPath   = "logPath"
	LoggerOptionOverwrite = "overwrite"
)

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
func (e *Logger) SaveOptions() []Option {
	return []Option{
		Option{
			Name:        LoggerOptionLogPath,
			Description: "full path to save the log",
			Type:        "string",
			Required:    true,
			Default:     ""},
		Option{
			Name:        LoggerOptionOverwrite,
			Description: "whether we should overwrite the log file if it exists or append",
			Type:        "bool",
			Required:    false,
			Default:     "false"}}
}

// ConfigureBackup configures the engine so it is ready to use
func (e *Logger) ConfigureSave(options map[string]interface{}) error {
	for k, v := range options {
		switch strings.ToLower(k) {
		case strings.ToLower(LoggerOptionLogPath):
			vString, ok := v.(string)
			if !ok {
				return goblerr.New("Invalid option", ErrorInvalidOptionValue, nil, LoggerOptionLogPath+" must be a string")
			}
			e.logPath = vString
		case strings.ToLower(LoggerOptionOverwrite):
			vbool, ok := v.(bool)
			if !ok {
				return goblerr.New("Invalid option", ErrorInvalidOptionValue, nil, LoggerOptionOverwrite+" must be a bool")
			}
			e.overWrite = vbool
		}
	}

	if e.logPath == "" {
		return goblerr.New("Must provide log path", ErrorRequiredOptionMissing, nil, LoggerOptionLogPath+" is required")
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
func (e *Logger) ShouldSave(fileSig files.Signature) (bool, error) {
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
