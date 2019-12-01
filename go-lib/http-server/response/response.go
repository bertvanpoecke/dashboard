package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type jsonError struct {
	Error errorDetails `json:"error"`
}

type errorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteJSONError writes an elegant JSON-Formatted Error to the given responsewriter
func WriteJSONError(w http.ResponseWriter, errorCode int, errorMessage string) {
	errDetails := errorDetails{
		Message: errorMessage,
		Code:    fmt.Sprint(errorCode),
	}

	errOut := jsonError{
		Error: errDetails,
	}

	WriteJSON(w, errorCode, errOut)
}

// WriteJSON encodes the given output-struct to json and writes it to the given responsewriter
func WriteJSON(w http.ResponseWriter, httpStatusCode int, object interface{}) {
	data, err := json.Marshal(object)
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(httpStatusCode)
		_, err = w.Write(data)
		if err != nil {
			logrus.Errorf("Failed writing json-bytes to responseWriter: %v", err)
		}
	} else {
		logrus.Errorf("Failed marshalling object to JSON: %v", err)
		w.WriteHeader(500)
	}
}

// WriteFile writes bytes to a file that is then disposed by the given responsewriter
func WriteFile(w http.ResponseWriter, fileBytes []byte, fileType, filename string) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))
	w.Header().Set("Content-Type", fileType)
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	_, err := w.Write(fileBytes)
	if err != nil {
		logrus.Errorf("Failed writing file: %v", err)
	}
}

// WriteJSONFile writes json content to a file that is then disposed by the given responsewriter
func WriteJSONFile(w http.ResponseWriter, json []byte, filename string) {
	WriteFile(w, json, "application/json", filename)
}

// WriteXLSXFile writes xlsx content to a file that is then disposed by the given responsewriter
func WriteXLSXFile(w http.ResponseWriter, xlsx []byte, filename string) {
	WriteFile(w, xlsx, "application/vnd.ms-excel", filename)
}

// WritePDFFile writes pdf content to a file that is then disposed by the given responsewriter
func WritePDFFile(w http.ResponseWriter, pdf []byte, filename string) {
	WriteFile(w, pdf, "application/pdf", filename)
}

// WriteCSVFile writes csv content to a file that is then disposed by the given responsewriterr
func WriteCSVFile(w http.ResponseWriter, csv []byte, filename string) {
	WriteFile(w, csv, "text/csv", filename)
}

// WriteHTMLFile writes html content to a file that is then disposed by the given responsewriter
func WriteHTMLFile(w http.ResponseWriter, html []byte, filename string) {
	WriteFile(w, html, "text/html", filename)
}

// WriteXMLFile writes xml content to a file that is then disposed by the given responsewriter
func WriteXMLFile(w http.ResponseWriter, xml []byte, filename string) {
	WriteFile(w, xml, "application/xml", filename)
}
