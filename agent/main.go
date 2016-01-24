package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sethjback/gobble/agent/apihandler"
	"github.com/sethjback/gobble/agent/manager"
	"github.com/sethjback/gobble/httpapi"
	"github.com/sethjback/gobble/spec"
	"github.com/sethjback/gobble/util/log"
	"github.com/sethjback/gobble/version"
)

type config struct {
	IP             string           `toml:"ip"`
	PORT           string           `toml:"port"`
	PrivateKeyPath string           `toml:"privatekey"`
	Coordinator    spec.Coordinator `toml:"coordinator"`
	LogConfig      log.Config       `toml:"logging"`
}

func main() {

	var cPath string

	flag.StringVar(&cPath, "config", "", "Path to the config file")
	flag.Parse()

	conf, err := parseConfig(cPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Init(conf.LogConfig)
	log.Infof("main", "agent starting. Version: %s", version.Version.String())
	log.Debug("main", "config:", *conf)

	err = manager.Init(conf.PrivateKeyPath, conf.Coordinator)
	if err != nil {
		log.Fatalf("main", "Error initializing manager: %v", err)
	}

	httpAPI := new(httpapi.Server)

	httpAPI.Configure(apihandler.Routes)
	httpAPI.Start(conf.IP, conf.PORT)
}

func parseConfig(path string) (*config, error) {
	if path == "" {
		return nil, errors.New("Config path empty")
	}

	cFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf config

	if _, err := toml.Decode(string(cFile), &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
