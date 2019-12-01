package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

const ui = "./ui-dist/"

// optionsHandler to satisfy preflight requests
func (api *API) optionsHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Allow", "GET,POST,OPTIONS")
}

// index returns the api name and supported calls
func (api *API) index(w http.ResponseWriter, r *http.Request) {
	indexFile := ui + "index.html"
	body, err := ioutil.ReadFile(indexFile)
	if err != nil {
		logrus.Errorf("Failed to write index to responsewriter (%v)", err)
	}
	fmt.Fprint(w, string(body))
}

// WriteError writes an error to the given ResponseWriter
func WriteError(w http.ResponseWriter, code int, message string, args ...interface{}) {
	logrus.Errorf(message, args...)
	w.Header().Set("content-type", "application/json")
	http.Error(w, fmt.Sprintf(message, args...), code)
}
