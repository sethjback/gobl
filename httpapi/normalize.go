package httpapi

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/sethjback/gobl/goblerr"
)

const (
	ErrorDateRequired = "DateHeaderRequired"
	ErrorDateInvalid  = "DateHeaderInvalid"
)

// Normalize
// Middleware:
// It implements the ServeHTTP interface for negroni middleware. The function is to
// evaluate the incoming request and validate/store the required headers
type Normalize struct {
}

func NewNormalize() *Normalize {
	return &Normalize{}
}

// ServeHTTP is the interface implementation for negroni middleware
func (n Normalize) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	_, err := validateTimestamp(r.Header.Get(HeaderGoblDate))
	if err != nil {
		resp := Response{
			HTTPCode: 400,
			Error:    err,
		}

		resp.Write(rw)
		return
	}

	req := &Request{}
	req.Headers = r.Header

	if r.Body != nil {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

		if err != nil {
			resp := Response{

				HTTPCode: 413,
				Error:    goblerr.New("Body invalid", ErrorRequestBodyInvalid, "normalize", err),
			}
			resp.Write(rw)
			return
		}

		if err := r.Body.Close(); err != nil {
			resp := Response{
				HTTPCode: 400,
				Error:    goblerr.New("Error reading body", ErrorRequestBodyInvalid, "normalize", err),
			}
			resp.Write(rw)
			return
		}

		if len(body) != 0 {
			req.Body = bytes.NewReader(body)
		}
	}

	req.Query = r.URL.Query()
	req.Host = r.Host
	req.Path = r.URL.Path

	ctx := r.Context()
	next(rw, r.WithContext(context.WithValue(ctx, request, req)))
}

func validateTimestamp(timestamp string) (int, goblerr.Error) {
	if timestamp == "" {
		return -1, goblerr.New("Date header not set", ErrorDateRequired, "normalize", "you must provide the x-gobl-date header in every request")
	}

	tint, err := strconv.Atoi(timestamp)
	if err != nil {
		return -1, goblerr.New("Date header invalid", ErrorDateInvalid, "normalize", "header not a valid unix timestamp")
	}

	cTime := int(time.Now().UTC().Unix())
	if tint > cTime+5 {
		return -1, goblerr.New("Date header invalid", ErrorDateInvalid, "normalize", "timestamp cannot be in the future")
	}

	if tint < cTime-500 {
		return -1, goblerr.New("Date header invalid", ErrorDateInvalid, "normalize", "timestamp must be within 5 min of now")
	}

	return tint, nil
}
