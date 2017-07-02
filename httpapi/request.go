package httpapi

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/goblerr"
)

const (
	ErrorRequestBodyInvalid = "RequestBodyInvalid"
	ErrorRequestInvalid     = "InvalidRequest"
	ErrorRequestFailed      = "RequestFailed"

	HeaderGoblDate = "x-gobl-date"
	HeaderGoblSig  = "x-gobl-signature"
)

// requestParms is the key for standardized request parameters
const requestKey key = 1

// Request provides a starndardized way of accessing incoming requests
// We standardize information both for receiving and sending.
//
// Context:
// It implements context functions for the context.Context package so that
// the rest of gobl can access the standardized parameters via the http.Request context
type Request struct {
	Headers         http.Header
	RouteParameters httprouter.Params
	Body            io.ReadSeeker
	Host            string
	Path            string
	Method          string
	Query           url.Values
	Client          *http.Client
}

func NewRequest(host, path, method string) *Request {
	return &Request{
		Headers: http.Header{},
		Host:    host,
		Path:    path,
		Method:  method,
	}
}

// RequestFromContext returns the reqest that has been stored in a context
func RequestFromContext(ctx context.Context) *Request {
	return ctx.Value(requestKey).(*Request)
}

func (r *Request) JsonBody(decodeTo interface{}) error {
	if r.Body == nil {
		return nil
	}

	if cType := r.Headers.Get("Content-Type"); cType != "application/json" {
		return goblerr.New("Body must be valid json", ErrorRequestBodyInvalid, "Content-Type must be application/json")
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return goblerr.New("Body not valid json", ErrorRequestBodyInvalid, "read body failed")
	}

	if err = json.Unmarshal(b, decodeTo); err != nil {
		return goblerr.New("Body not valid json", ErrorRequestBodyInvalid, "json unmarshal failed: "+err.Error())
	}

	return nil
}
