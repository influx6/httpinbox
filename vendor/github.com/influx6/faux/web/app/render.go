package app

import (
	"encoding/json"
	"io"
	"net/http"
)

// Field defines a struct for collating fields errors that occur.
type Field struct {
	Name     string      `json:"field_name"`
	Value    string      `json:"field_value"`
	Error    string      `json:"field_error"`
	Expected interface{} `json:"expected_value"`
}

// JSONError defines a json error response struct
type JSONError struct {
	Error  string  `json:"error"`
	Fields []Field `json:"fields,omitempty"`
}

// Render writes the giving data into the response as JSON.
func Render(code int, r *http.Request, w http.ResponseWriter, data interface{}) {
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	jsd, err := json.Marshal(data)
	if err != nil {
		data = []byte("{}")
	}

	if cb := r.URL.Query().Get("callback"); cb != "" {
		io.WriteString(w, cb+"("+string(jsd)+")")
		return
	}

	io.WriteString(w, string(jsd))
}

// RenderResponse writes the giving data into the response as JSON to the passed
// ResponseRequest.
func RenderResponse(code int, r *ResponseRequest, data interface{}) {
	Render(code, r.R, r, data)
}

// RenderResponseErrorWithStatus renders the giving error as a json response to
// the ResponseRequest.
func RenderResponseErrorWithStatus(status int, err error, r *ResponseRequest) {
	Render(status, r.R, r, JSONError{Error: err.Error()})
}

// RenderErrorWithStatus renders the giving error as a json response.
func RenderErrorWithStatus(status int, err error, r *http.Request, w http.ResponseWriter) {
	Render(status, r, w, JSONError{Error: err.Error()})
}

// RenderError renders the giving error as a json response.
func RenderError(err error, r *http.Request, w http.ResponseWriter) {
	Render(http.StatusBadRequest, r, w, JSONError{Error: err.Error()})
}

// RenderResponseError renders the giving error as a json response to the
// passed ResponseRequest object.
func RenderResponseError(err error, r *ResponseRequest) {
	Render(http.StatusBadRequest, r.R, r, JSONError{Error: err.Error()})
}
