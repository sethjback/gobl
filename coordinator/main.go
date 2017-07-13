package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/coordinator/apihandler"
	"github.com/sethjback/gobl/coordinator/cron"
	"github.com/sethjback/gobl/coordinator/db"
	"github.com/sethjback/gobl/coordinator/email"
	"github.com/sethjback/gobl/coordinator/grpcserver"
	"github.com/sethjback/gobl/httpapi"
)

func main() {

	var configPath string
	flag.StringVar(&configPath, "config", "", "path to the config file")
	flag.Parse()

	cMap, err := config.Map(configPath)
	if err != nil {
		fmt.Printf("error with config: %v\n", err)
		os.Exit(1)
	}

	cs := config.New()

	err = db.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring DB: %v\n", err)
		os.Exit(1)
	}

	err = email.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring email: %v\n", err)
		os.Exit(1)
	}

	err = httpapi.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring api server: %v\n", err)
		os.Exit(1)
	}

	err = grpcserver.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring grpc server: %v\n", err)
		os.Exit(1)
	}

	err = db.Init(cs)
	if err != nil {
		fmt.Printf("Error db init: %v\n", err)
		os.Exit(1)
	}
	grpcserver.Init(db.Get())
	cron.Init(db.Get())
	apihandler.Init(db.Get())

	httpAPI := httpapi.New(apihandler.Routes)
	httpAPI.Start(cs, func() {
		fmt.Println("Shutting Down")
		grpcserver.Shutdown()
		cron.Shutdown()
		db.Shutdown()
	})
}
