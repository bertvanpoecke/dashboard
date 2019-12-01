package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger wraps a log around a handler when called
func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		logrus.Infof("%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
