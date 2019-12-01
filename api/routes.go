package api

import (
	"net/http"

	"github.com/bertvanpoecke/dashboard/go-lib/http-server/router"
)

// getRoutes holds the http_server routes definition
func (api *API) getRoutes() []router.Route {
	return []router.Route{
		{
			Method:      http.MethodGet,
			Path:        "/",
			Description: "Index",
			Protected:   false,
			Protection:  []string{"scope"},
			HandlerFunc: api.index,
		},
	}
}
