package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/sethjback/gobl/agent/apihandler"
	"github.com/sethjback/gobl/agent/coordinator"
	"github.com/sethjback/gobl/agent/grpc"
	"github.com/sethjback/gobl/agent/manager"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/email"
	"github.com/sethjback/gobl/goblgrpc"
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

	cs := config.SetEnvMap(cMap, config.New())
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

	err = grpc.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring grpc server: %v\n", err)
		os.Exit(1)
	}

	err = coordinator.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring coordinator: %v\n", err)
		os.Exit(1)
	}

	err = manager.Init(cs)
	if err != nil {
		fmt.Printf("Error initializing manager: %v\n", err)
		os.Exit(1)
	}

	gs, err := grpc.New(cs)
	if err != nil {
		fmt.Printf("Error starting grpc server: %v\n", err)
		os.Exit(1)
	}

	go gs.StartServer()

	// just for testing
	con, err := gs.Client(coordinator.FromConfig(cs))
	if err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}
	msg, err := con.Finish(context.Background(), &goblgrpc.JobDefinition{Id: "asdf"})
	fmt.Printf("msg: %+v\n", msg)

	// it works
	httpAPI := httpapi.New(apihandler.Routes)
	httpAPI.Start(cs, func() {
		fmt.Println("Shutting Down")
		manager.Shutdown()
		gs.StopServer()
	})
}
