package httpapi

// Route defines a single endpoint
type Route struct {
	// Method is the HTTP method to listen on
	Method string

	// Path indictes the URL path. Wildcards in the form of :segment are supported
	// For more information see:
	// https://github.com/julienschmidt/httprouter
	Path string

	// Handler to use for this route
	Handler RouteHandler
}

// RouteHandler is the definition functions must meet to handle incoming requests
type RouteHandler func(*Request) Response
