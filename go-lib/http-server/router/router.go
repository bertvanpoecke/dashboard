package router

import (
	"fmt"
	"net/http"

	"github.com/bertvanpoecke/dashboard/go-lib/http-server/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Route is a type that is used to register endpoints in a Router
type Route struct {
	Name        string
	Method      string
	Path        string
	Description string
	Protected   bool
	Protection  []string
	Log         bool
	HandlerFunc http.HandlerFunc
}

// New returns a new mux HTTP Router to handle the given routes/paths.
// if no valid optionsHandler is given an error is returned
func New(routes []Route, optionsHandler http.HandlerFunc) (*mux.Router, error) {
	// check if given optionsHandler is not nil
	if optionsHandler == nil {
		return nil, fmt.Errorf("OptionsHandler should not be nil")
	}

	// Init router
	router := mux.NewRouter().StrictSlash(true)

	// Register preflight requests handler
	var options http.Handler
	options = http.HandlerFunc(optionsHandler)
	options = middleware.CORS(options)
	router.
		Methods(http.MethodOptions).
		Handler(options)

	// Register all custom routes
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		// if logger is required, add middleware
		if route.Log {
			handler = middleware.Logger(handler)
		}

		// Add general middleware
		handler = middleware.CORS(handler)
		handler = handlers.CompressHandler(handler)

		// Add the route to the router
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(handler)

		logrus.Infof("Registered %v %v on the router", route.Method, route.Path)
	}
	return router, nil
}
