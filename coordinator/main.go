package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/fatih/structs"
	"github.com/sethjback/gobl/coordinator/apihandler"
	"github.com/sethjback/gobl/coordinator/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/util/log"
	"github.com/sethjback/gobl/version"
)

type config struct {
	IP             string     `toml:"ip"`
	PORT           string     `toml:"port"`
	DBPath         string     `toml:"dbpath"`
	PrivateKeyPath string     `toml:"privatekey"`
	LogConfig      log.Config `toml:"logging"`
}

func main() {

	var cPath string

	flag.StringVar(&cPath, "config", "", "Path to the config file")
	flag.Parse()

	conf, err := parseConfig(cPath)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		os.Exit(1)
	}

	log.Init(conf.LogConfig)
	log.Infof("main", "coordinator starting. Version: %s", version.Version.String())
	log.Debug("main", "config:", *conf)

	err = manager.Init(structs.Map(conf))
	if err != nil {
		log.Fatal("main", "Error initializing manager: %v", err)
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
