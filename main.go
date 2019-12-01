package main

import (
	"syscall"

	"github.com/bertvanpoecke/dashboard/api"
	"github.com/bertvanpoecke/dashboard/go-lib/shutdown"
	"github.com/sirupsen/logrus"
)

// Version is filled in by Jenkins when building
var version string

// global log
var logger = logrus.New()

// configure death: pass the signals you want to end your application
var httpServerShutdown = shutdown.NewShutdown(syscall.SIGINT, syscall.SIGTERM)

func main() {
	httpServerShutdown.SetLogger(logger)
	logger.Infof("API - Version %v", version)

	// Init loads config and inits new API
	localAPI, err := api.InitNewAPI()
	if err != nil {
		logger.Fatal(err)
	}

	// start api
	go func(localAPI *api.API) {
		if err := localAPI.Start(); err != nil {
			logger.Errorf("failed starting API: %v", err)
			httpServerShutdown.Start()
		}
	}(localAPI)

	// wait for shutdown
	httpServerShutdown.WaitForShutdownWithFunc(func() {
		if localAPI != nil {
			if err := localAPI.Stop(); err != nil {
				logger.Errorf("failed stopping API: %v", err)
			}
		}
	})
}
