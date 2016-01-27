package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sethjback/gobl/util/log"
)

// Server is the main handler for the api
type Server struct {
	router *mux.Router
}

// RouteHandler to hand api functions
type RouteHandler func(http.ResponseWriter, *http.Request) (*APIResponse, error)

// APIResponse is the common format used for the api response
type APIResponse struct {
	Message  string                 `json:"message"`
	HTTPCode int                    `json:"-"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

func (e *APIResponse) Error() string {
	return e.Message
}

// Route contains information for a single api endpoint
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc RouteHandler
}

// Routes is a slice of Route
type Routes []Route

// NewError returns a new response struct with the source error saved as data payload
func NewError(source string, message string, code int) *APIResponse {
	return &APIResponse{
		Data:     map[string]interface{}{"source": source},
		Message:  message,
		HTTPCode: code}
}

//Configure sets up the router
func (s *Server) Configure(routes Routes) {

	s.router = mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		s.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(apiHandler(route.HandlerFunc))
	}

}

func apiHandler(f RouteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := f(w, r)

		w.Header().Set("Cache-Control", "must-revalidate")
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			if e, ok := err.(*APIResponse); ok {
				errString, _ := json.Marshal(err)
				w.WriteHeader(e.HTTPCode)
				w.Write(errString)
			} else {
				w.WriteHeader(500)
				w.Write([]byte(`{"message": "Server Error"}`))
			}
		} else {

			if len(resp.Message) == 0 {
				resp.Message = "success"
			}

			j, err := json.Marshal(resp)
			if err != nil {
				log.Debug("apiHandler", err.Error())
				j = []byte(`{"message": "Trouble marshalling the success response"}`)
			}
			w.WriteHeader(resp.HTTPCode)
			w.Write(j)
		}
	}
}

// Start serving the api
func (s *Server) Start(ip string, port string) {
	var addr string
	if len(ip) == 0 {
		addr = "" + ":" + port
	} else {
		addr = ip + ":" + port
	}

	http.ListenAndServe(addr, s.router)
}
