// Package log  contains convienience methods for logging pertinant information
// Kudos to Pual Ruane's TMSU for the general outline of this package
// github.com/oniony/TMSU
package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sethjback/gobl/config"
)

type level struct {
	Fatal int
	Error int
	Warn  int
	Info  int
	Debug int
}

var conf *config.Log
var out io.Writer

// Level is the log level to use
var Level = level{1, 2, 3, 4, 5}

// Init intilize the logging framework
// Output will ultimately determine where to log. For now, everything logs to Stdout
// level determins the level of logging
// verbosity determins how much information to put in the log (right now it just controls timestamp display)
func Init(config config.Log) {
	conf = &config
	out = os.Stdout
}

// Fatal will log one last mesasge, then exit with an error code
func Fatal(source string, values ...interface{}) {
	if conf.Level >= Level.Fatal {
		log(source, values)
	}
	os.Exit(1)
}

// Fatalf will log one last formated message then exit with error code
func Fatalf(source string, format string, values ...interface{}) {
	if conf.Level >= Level.Fatal {
		logf(source, format, values)
	}
	os.Exit(1)
}

// Error will log error level
func Error(source string, values ...interface{}) {
	if conf.Level >= Level.Error {
		log(source, values...)
	}
}

// Errorf will log formated error messages
func Errorf(source string, format string, values ...interface{}) {
	if conf.Level >= Level.Error {
		logf(source, format, values)
	}
}

// Warn will log warnings
func Warn(source string, values ...interface{}) {
	if conf.Level >= Level.Warn {
		log(source, values...)
	}
}

// Warnf will log formated warnings
func Warnf(source string, format string, values ...interface{}) {
	if conf.Level >= Level.Warn {
		logf(source, format, values)
	}
}

// Info will log informational messages
func Info(source string, values ...interface{}) {
	if conf.Level >= Level.Info {
		log(source, values...)
	}
}

// Infof will log formated informational messages
func Infof(source string, format string, values ...interface{}) {
	if conf.Level >= Level.Info {
		logf(source, format, values)
	}
}

// Debug will log messages at the debug leve
func Debug(source string, values ...interface{}) {
	if conf.Level >= Level.Debug {
		log(source, values...)
	}
}

// Debugf will log formated messages at the debug level
func Debugf(source string, format string, values ...interface{}) {
	if conf.Level >= Level.Debug {
		logf(source, format, values)
	}
}

func log(source string, values ...interface{}) {
	st := ""
	if conf.Verbosity > 1 {
		st = time.Now().String() + " "
	}

	st += "gobl: (" + source + ") " + strings.Repeat("%v ", len(values)) + "\n"

	fmt.Fprintf(out, st, values...)
}

func logf(source string, format string, values ...interface{}) {
	st := ""
	if conf.Verbosity > 1 {
		st = time.Now().String() + " "
	}

	st += "gobl: (" + source + ") " + format + "\n"
	fmt.Fprintf(out, st, values...)
}
