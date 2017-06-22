package httpapi

import (
	"net/http"
	"time"

	"github.com/sethjback/gobl/config"

	"github.com/julienschmidt/httprouter"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/tylerb/graceful"
	"github.com/urfave/negroni"
)

type key int

const serverConfigKey key = 0

type serverConfig struct {
	listen string
}

func SaveConfig(cs config.Store, env map[string]string) error {
	for k, v := range env {
		if k == "API_LISTEN" {
			sc := &serverConfig{}
			sc.listen = v
			cs.Add(serverConfigKey, sc)
			break
		}
	}

	return nil
}

func configFromStore(cs config.Store) *serverConfig {
	if sc, ok := cs.Get(serverConfigKey); ok {
		return sc.(*serverConfig)
	}
	return nil
}

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
func (s *Server) Start(cs config.Store, shutdown func()) {
	sc := configFromStore(cs)
	var listen string
	if sc != nil {
		listen = sc.listen
	} else {
		listen = ":8080"
	}

	n := negroni.New()

	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(NewNormalize())
	n.UseHandler(s.router)

	svr := graceful.Server{
		Timeout:           time.Duration(2) * time.Second,
		ShutdownInitiated: shutdown,
		Server:            &http.Server{Addr: listen, Handler: n},
	}

	svr.ListenAndServe()
}

// wrapRoute returns a httprouter appropriate handler
func wrapRoute(rh RouteHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req := RequestFromContext(r.Context())
		resp := rh(req, ps)
		resp.Write(w)
	}
}

func corsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setHeaders(w)
	w.WriteHeader(http.StatusOK)
}
