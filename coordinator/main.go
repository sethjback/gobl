package main

import (
	"flag"
	"fmt"

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
		fmt.Printf("error with config: %v", err)
	}

	cs := config.New()

	err = gobldb.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring DB: %v", err)
	}

	err = email.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring email: %v", err)
	}

	err = httpapi.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring api server: %v", err)
	}

	err = grpc.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring grpc server: %v", err)
	}

	err = manager.Init(cs)
	if err != nil {
		fmt.Printf("Error initializing manager: %v", err)
	}

	gs, err := grpc.New(cs)
	if err != nil {
		fmt.Printf("Error starting grpc server: %v", err)
	}

	go gs.StartServer()

	httpAPI := httpapi.New(apihandler.Routes)
	httpAPI.Start(cs, func() {
		fmt.Println("Shutting Down")
		manager.Shutdown()
		gs.StopServer()
	})
}
