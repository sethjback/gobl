package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/sethjback/gobl/agent/grpcclient"
	"github.com/sethjback/gobl/agent/grpcserver"
	"github.com/sethjback/gobl/agent/job"
	"github.com/sethjback/gobl/config"
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

	err = grpcserver.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring grpc server: %v\n", err)
		os.Exit(1)
	}

	err = grpcclient.SaveConfig(cs, cMap)
	if err != nil {
		fmt.Printf("Error configuring grpc client: %v\n", err)
		os.Exit(1)
	}

	job.Init()
	err = grpcserver.Init()
	if err != nil {
		fmt.Printf("Error starting grpc server: %s\n", err)
		os.Exit(1)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals

	job.Shutdown()
	grpcserver.Shutdown()

	os.Exit(0)
}
