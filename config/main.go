package config

import (
	"errors"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Config contains all the possible config options
type Config struct {
	Server      Server      `toml:"server"`
	DB          DB          `toml:"db"`
	Log         Log         `toml:"logging"`
	Email       Email       `toml:"email"`
	Coordinator Coordinator `toml:"coordinator"`
}

// Server config
type Server struct {

	// Listen defines what address:port to listen on
	Listen string `toml:"listen"`

	// Compress the server output
	Compress bool `toml:"compress"`

	// ShutdownWait is the number of seconds to wait for backup processes
	// to finish before stopping the server
	ShutdownWait int `toml:"shutdown_wait"`
}

type Agent struct {
	// PrivateKey is used to when communicating with the coordinator
	PrivateKey string `toml:"private_key"`
}

// Coordinator config.
// This is used to tell gobl agents about the coordinator so that they can
// confirm its identity when it connects
type Coordinator struct {

	// Address defines where the connections will be comming from
	Address string `toml:"address"`

	// PublicKey is the path to the coordinator's public key
	PublicKey string `toml:"public_key"`
}

// DB Config
type DB struct {

	// Path to the database file
	Path string `toml:"path"`
}

// Log defines the logging paramiters
type Log struct {
	Level     int    `toml:"level"`
	Verbosity int    `toml:"verbosity"`
	Output    string `toml:"output"`
}

// Email config
type Email struct {
	Server         string `toml:"server"`
	From           string `toml:"from"`
	To             string `toml:"to"`
	Subject        string `toml:"subject"`
	Authentication bool   `toml:"auth"`
	User           string `toml:"user"`
	Password       string `toml:"password"`
}

// Configured returns true if we have enough information to attempt to send emails
func (e *Email) Configured() bool {
	return len(e.Server) != 0 && len(e.From) != 0 && len(e.To) != 0
}

// ParseConfig parses a config file and returns a config object
func Parse(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("Config path empty")
	}

	cFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config

	if _, err := toml.Decode(string(cFile), &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
