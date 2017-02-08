package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/sethjback/gobl/goblerr"
)

// Response is the standardized response format sent from the API
type Response struct {
	// Data is API payload
	Data map[string]interface{} `json:"data"`

	// Error holds any errors that were encountered processing the request
	Error goblerr.Error `json:"error,omitempty"`

	// HTTPCode allows individual handlers to indiciate the appropriate http status code to return
	HTTPCode int `json:"-"`
}

// Write the standardized response to the given response writer
func (r *Response) Write(rw http.ResponseWriter) {
	setHeaders(rw)

	var j []byte
	var jErr error
	if r.Error != nil {
		j, jErr = json.Marshal(r.Error)
	} else {
		j, jErr = json.Marshal(r.Data)
	}
	if jErr != nil {
		j = []byte(`{"message": "Trouble marshalling success response"}`)
		r.HTTPCode = 500
	}
	rw.WriteHeader(r.HTTPCode)
	rw.Write(j)
}

func setHeaders(w http.ResponseWriter) {
	allowMethod := "GET, POST, OPTIONS"
	allowHeaders := "Content-Type, Authorization, x-gobl-date, x-gobl-signature"
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Allow", allowMethod)
	w.Header().Set("Access-Control-Allow-Methods", allowMethod)
	w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}
