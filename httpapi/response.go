package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/sethjback/gobl/goblerr"
)

// Response is the standardized response format sent from the API
type Response struct {
	// Message contains a generic message if the request is successful
	// and a user friendly error message if there was an error
	Message string `json:"message"`

	// Data is API payload
	Data map[string]interface{} `json:"data"`

	// Error holds any errors that were encountered processing the request
	Error goblerr.Error `json:"errors"`

	// HTTPCode allows individual handlers to indiciate the appropriate http status code to return
	HTTPCode int `json:"-"`
}

// Write the standardized response to the given response writer
func (r *Response) Write(rw http.ResponseWriter) {
	setHeaders(rw)

	j, jErr := json.Marshal(r)
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
