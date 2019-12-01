package api

import (
	"fmt"

	"net/http"

	"github.com/bertvanpoecke/dashboard/go-lib/http-server/router"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// API is a type to serve http content according to a config
type API struct {
	config *config
	router *mux.Router
	server *http.Server
}

// InitNewAPI loads configurations and inits the API accordingly
func InitNewAPI() (*API, error) {
	// load configuration
	config, err := loadConfiguration()
	if err != nil {
		return nil, err
	}
	// parse loglevel
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, err
	}
	logrus.SetLevel(logLevel)

	// create new API
	api := &API{
		config: config,
	}

	// create router
	router, err := router.New(api.getRoutes(), api.optionsHandler)
	if err != nil {
		return nil, err
	}

	// install fileserver
	fileserver := http.FileServer(http.Dir("./ui-dist/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileserver))

	api.router = router

	// start traffic getter
	traffic := NewTrafficInfoGetter()
	traffic.Start()

	// create server
	api.server = &http.Server{}
	api.server.Addr = fmt.Sprintf(":%v", api.config.Port)
	api.server.Handler = api.router

	// return our API!
	return api, nil
}

// Start the http_server ListenAndServe method for our API.
// This call is blocking.
func (api *API) Start() error {
	return api.server.ListenAndServe()
}

// Stop shuts down the API
func (api *API) Stop() error {
	return api.server.Shutdown(nil)
}
