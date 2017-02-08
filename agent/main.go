package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sethjback/gobl/agent/apihandler"
	"github.com/sethjback/gobl/agent/manager"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/util/log"
	"github.com/sethjback/gobl/version"
)

func main() {

	var cPath string

	flag.StringVar(&cPath, "config", "", "Path to the config file")
	flag.Parse()

	conf, err := config.Parse(cPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Init(conf.Log)
	log.Infof("main", "agent starting. Version: %s", version.Version.String())
	log.Debug("main", "config:", *conf)

	err = manager.Init(conf)
	if err != nil {
		log.Fatalf("main", "Error initializing manager: %v", err)
	}

	httpAPI := httpapi.New(apihandler.Routes)

	httpAPI.Start(conf.Server, func() {
		log.Infof("main", "shutting down")
		manager.Shutdown()
	})
}
