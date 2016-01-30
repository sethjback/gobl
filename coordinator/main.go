package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/coordinator/apihandler"
	"github.com/sethjback/gobl/coordinator/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/util/log"
	"github.com/sethjback/gobl/version"
)

func main() {

	var cPath string

	flag.StringVar(&cPath, "config", "", "Path to the config file")
	flag.Parse()

	conf, err := config.ParseConfig(cPath)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		os.Exit(1)
	}

	log.Init(conf.Log)
	log.Infof("main", "coordinator starting. Version: %s", version.Version.String())
	log.Debug("main", "config:", *conf)

	err = manager.Init(conf)
	if err != nil {
		log.Fatalf("main", "Error initializing manager: %v", err)
	}

	httpAPI := new(httpapi.Server)

	address := strings.Split(conf.Host.Address, ":")
	if len(address) != 2 {
		log.Fatalf("main", "Invalid host address. Must be in form ip:port")
	}

	httpAPI.Configure(apihandler.Routes)
	httpAPI.Start(address[0], address[1])

}
