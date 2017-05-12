package httpapi

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/sethjback/gobl/config"
	"github.com/tylerb/graceful"
	"github.com/urfave/negroni"
)

// Server handles accepting and replying to API requests
type Server struct {
	router *httprouter.Router
}

// New returns a new httpapi.Server configured to respond to the routes provided
func New(routes []Route) *Server {
	s := &Server{router: httprouter.New()}

	for _, r := range routes {
		switch r.Method {
		case "GET":
			s.router.GET(r.Path, wrapRoute(r.Handler))
		case "POST":
			s.router.POST(r.Path, wrapRoute(r.Handler))
		case "PUT":
			s.router.PUT(r.Path, wrapRoute(r.Handler))
		case "DELETE":
			s.router.DELETE(r.Path, wrapRoute(r.Handler))
		}
	}

	s.router.OPTIONS("/*all", corsHandler)

	return s
}

// Start listening on given address.
// The shutdown function will be called before the server exits
func (s *Server) Start(c config.Server, shutdown func()) {
	n := negroni.New()

	if c.Compress {
		n.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	n.Use(NewNormalize())
	n.UseHandler(s.router)

	graceful.Run(c.Listen, time.Duration(c.ShutdownWait)*time.Second, n)
}

// wrapRoute returns a httprouter appropriate handler
func wrapRoute(rh RouteHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		resp := rh(RequestFromContext(r.Context()), ps)
		resp.Write(w)
	}
}

func corsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setHeaders(w)
	w.WriteHeader(http.StatusOK)
}
