package middleware

import "net/http"

const (
	corsAllowHeaders     = "authorization, x-bemobile-language, x-bemobile-guid, x-bemobile-source, x-bemobile-workspaceid, x-bemobile-workspaces, x-bemobile-ownerid, content-type, content-language, accept-language"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

/* The CORS handler provides the Access Control headers on the HTTP methods configured in the router.
You will need a separate OPTIONS handler if the OPTIONS method isn't part of the configured methods on a route.
*/

// CORS middleware
func CORS(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		writer.Header().Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		writer.Header().Set("Access-Control-Allow-Headers", corsAllowHeaders)
		writer.Header().Set("Access-Control-Allow-Methods", corsAllowMethods)
		writer.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)

		next.ServeHTTP(writer, req)
	})
}
