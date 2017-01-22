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

type LogLine struct {
	Start         string   `json:"start"`
	End           string   `json:"end"`
	File          string   `json:"file"`
	Modifications []string `json:"modifications"`
	Engines       []string `json:"engines"`
	Signature     string   `json:"signature"`
	Chunks        int      `json:"chunks"`
	Bytes         int      `json:"bytes"`
}

// Name return "Logger"
func (e *Logger) Name() string {
	return "Logger"
}

// BackupOptions lists all the availalbe options
func (e *Logger) SaveOptions() []Option {
	return getOptions()
}

// ConfigureBackup configures the engine so it is ready to use
func (e *Logger) ConfigureSave(options map[string]interface{}) error {
	return e.configure(options)
}

// ShouldBackup always returns true: we want to log each file
func (e *Logger) ShouldSave(file files.File) (bool, error) {
	return true, nil
}

// Backup collects information about the file then writes it to a log once the backup has finished
func (e *Logger) Save(reader io.Reader, file files.File, errc chan<- error) {
	e.recordAndSave(reader, file, errc)
}

// Retrieve does nothing
func (e *Logger) Retrieve(file files.File) (io.Reader, error) {
	return nil, errors.New("Cannot use Logger to restore files")
}

func (e *Logger) RestoreOptions() []Option {
	return getOptions()
}

func (e *Logger) ConfigureRestore(options map[string]interface{}) error {
	return e.configure(options)
}

func (e *Logger) ShouldRestore(file files.File) (bool, error) {
	return true, nil
}

func (e *Logger) Restore(reader io.Reader, file files.File, errc chan<- error) {
	e.recordAndSave(reader, file, errc)
}

func (e *Logger) recordAndSave(reader io.Reader, file files.File, errc chan<- error) {
	fileData := new(LogLine)
	fileData.File = file.Path
	fileData.Modifications = file.Signature.Modifications
	sig, _ := json.Marshal(file.Signature)
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
		fileData.Bytes += len(buf)
	}

	fileData.End = time.Now().String()

	dataBytes, _ := json.Marshal(fileData)

	lfile, err := os.OpenFile(e.logPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		errc <- err
	}
	_, err = lfile.WriteString(string(dataBytes) + "\n")
	lfile.Close()
}

func (e *Logger) configure(options map[string]interface{}) error {
	for k, v := range options {
		switch strings.ToLower(k) {
		case strings.ToLower(LoggerOptionLogPath):
			vString, ok := v.(string)
			if !ok {
				return goblerr.New("Invalid option", ErrorInvalidOptionValue, "logger", LoggerOptionLogPath+" must be a string")
			}
			e.logPath = vString
		case strings.ToLower(LoggerOptionOverwrite):
			vbool, ok := v.(bool)
			if !ok {
				return goblerr.New("Invalid option", ErrorInvalidOptionValue, "logger", LoggerOptionOverwrite+" must be a bool")
			}
			e.overWrite = vbool
		}
	}

	if e.logPath == "" {
		return goblerr.New("Must provide log path", ErrorRequiredOptionMissing, "logger", LoggerOptionLogPath+" is required")
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

func getOptions() []Option {
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
