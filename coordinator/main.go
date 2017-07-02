package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/coordinator/apihandler"
	"github.com/sethjback/gobl/coordinator/grpc"
	"github.com/sethjback/gobl/coordinator/manager"
	"github.com/sethjback/gobl/email"
	"github.com/sethjback/gobl/gobldb"
	"github.com/sethjback/gobl/httpapi"
)

func main() {

	var configPath string
	flag.StringVar(&configPath, "config", "", "path to the config file")
	flag.Parse()

	cMap, err := config.Map(configPath)
	if err != nil {
		fmt.Printf("error with config: %v\n", err)
	}

	cs := config.New()

	err = gobldb.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring DB: %v\n", err)
	}

	err = email.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring email: %v\n", err)
	}

	err = httpapi.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring api server: %v\n", err)
	}

	err = grpc.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring grpc server: %v\n", err)
		os.Exit(1)
	}

	err = manager.Init(cs)
	if err != nil {
		fmt.Printf("Error initializing manager: %v\n", err)
		os.Exit(1)
	}

	httpAPI := httpapi.New(apihandler.Routes)
	httpAPI.Start(cs, func() {
		fmt.Println("Shutting Down")
		manager.Shutdown()
	})
}
