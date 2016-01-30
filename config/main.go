package config

import (
	"errors"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Config contains all the possible config options
type Config struct {
	Host        Host        `toml:"host"`
	DB          DB          `toml:"db"`
	Log         Log         `toml:"logging"`
	Email       Email       `toml:"email"`
	Coordinator Coordinator `toml:"coordinator"`
}

// Host config
type Host struct {
	Address        string `toml:"address"`
	PrivateKeyPath string `toml:"privatekey"`
}

// Coordinator config
type Coordinator struct {
	Address       string `toml:"address"`
	PublicKeyPath string `toml:"publickey"`
}

// DB Config
type DB struct {
	DBPath string `toml:"dbpath"`
}

// Log defines the logging paramiters
type Log struct {
	Level     int    `toml:"level"`
	Verbosity int    `toml:"verbosity"`
	Output    string `toml:"output"`
}

// Email config
type Email struct {
	ServerAddress  string `toml:"server"`
	From           string `toml:"from"`
	To             string `toml:"to"`
	Subject        string `toml:"subject"`
	Authentication bool   `toml:"auth"`
	User           string `toml:"user"`
	Password       string `toml:"password"`
}

// Configured returns true if we have enough information to attempt to send emails
func (e *Email) Configured() bool {
	return len(e.ServerAddress) != 0 && len(e.From) != 0 && len(e.To) != 0
}

// ParseConfig parses a config file and returns a config object
func ParseConfig(path string) (*Config, error) {
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
