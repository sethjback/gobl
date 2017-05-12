package httpapi

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/util/log"
)

const (
	ErrorRequestBodyInvalid = "RequestBodyInvalid"
	ErrorRequestInvalid     = "InvalidRequest"
	ErrorRequestFailed      = "RequestFailed"

	HeaderGoblDate = "x-gobl-date"
	HeaderGoblSig  = "x-gobl-signature"
)

// paramsKey for use with context
type key int

// requestParms is the key for standardized request parameters
const request key = 1

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
	return ctx.Value(request).(*Request)
}

// String returns the request string that is appropriate for signing
func (r *Request) String() string {
	uri := strings.ToLower(r.Path)
	query := queryString(r.Query)

	var body string
	if r.Body != nil {
		body = bodyHash(r.Body)
	} else {
		body = ""
	}

	var headers bytes.Buffer
	headers.WriteString("authorization:" + r.Headers.Get("authorization") + "\n")
	headers.WriteString(HeaderGoblDate + ":" + r.Headers.Get(HeaderGoblDate))

	return strings.Join([]string{
		r.Method,
		r.Host,
		uri,
		query,
		headers.String(),
		body}, "\n")
}

func (r *Request) SetBody(body interface{}) error {
	bbytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	r.Body = bytes.NewReader(bbytes)

	return nil
}

func (r *Request) JsonBody(decodeTo interface{}) goblerr.Error {
	if cType := r.Headers.Get("Content-Type"); cType != "application/json" {
		return goblerr.New("Body must be valid json", ErrorRequestBodyInvalid, "request", "Content-Type must be application/json")
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return goblerr.New("Body not valid json", ErrorRequestBodyInvalid, "request", "read body failed")
	}

	if err = json.Unmarshal(b, decodeTo); err != nil {
		log.Debug(err.Error())
		return goblerr.New("Body not valid json", ErrorRequestBodyInvalid, "request", "json unmarshal failed")
	}

	return nil
}

// bodyHash returns the sha256 sum of the body
func bodyHash(reader io.ReadSeeker) string {
	hash := sha256.New()

	start, _ := reader.Seek(0, 1)
	defer reader.Seek(start, 0)

	io.Copy(hash, reader)
	s := hash.Sum(nil)
	return hex.EncodeToString(s)
}

// sorts the query parameters and sends them back in alpha order
func queryString(values url.Values) string {
	keys := make([]string, len(values))
	var buff bytes.Buffer
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	l := len(keys)
	for k, v := range keys {
		vals := values[v]
		sort.Strings(vals)
		l1 := len(vals)
		for k1, v1 := range vals {
			buff.WriteString(v + "=" + v1)
			if k1 != l1-1 {
				buff.WriteString("&")
			} else if k != l-1 {
				buff.WriteString("&")
			}
		}

	}

	return buff.String()
}

func (r *Request) Send(s keys.Signer) (*Response, goblerr.Error) {
	err := prepAndSign(r, s)
	if err != nil {
		return nil, goblerr.New("Unable to sign message", ErrorRequestFailed, "request", err)
	}
	switch r.Method {
	case "POST":
		return post(r)
	case "GET":
		return get(r)
	}
	return nil, goblerr.New("Invalid method", ErrorRequestInvalid, "request", "must be POST or GET")
}

func prepAndSign(r *Request, s keys.Signer) error {
	if d := r.Headers.Get(HeaderGoblDate); d == "" {
		r.Headers.Set(HeaderGoblDate, strconv.Itoa(int(time.Now().UTC().Unix())))
	}

	sig, err := s.Sign([]byte(r.String()))

	if err != nil {
		return err
	}

	r.Headers.Set(HeaderGoblSig, sig)

	return nil
}

// Post a request
func post(r *Request) (*Response, goblerr.Error) {
	req, err := http.NewRequest("POST", r.Host+r.Path, r.Body)
	if err != nil {
		return nil, goblerr.New("Invalid request", ErrorRequestInvalid, "request", err)
	}

	req.Header = r.Headers
	req.Header.Set(HeaderGoblDate, strconv.Itoa(int(time.Now().UTC().Unix())))
	req.Header.Set("Content-Type", "application/json")

	if r.Client == nil {
		r.Client = &http.Client{CheckRedirect: checkRedirect}
	} else {
		r.Client.CheckRedirect = checkRedirect
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, goblerr.New("Unable to send request", ErrorRequestFailed, "request", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, goblerr.New("Unable to read response", ErrorRequestFailed, "request", err)
	}

	resp.Body.Close()

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, goblerr.New("Unable to unmarshal", ErrorRequestFailed, "request", err)
	}
	response.HTTPCode = resp.StatusCode

	return &response, nil
}

// Get a request
func get(r *Request) (*Response, goblerr.Error) {
	req, err := http.NewRequest("GET", r.Host+r.Path, nil)
	if err != nil {
		return nil, goblerr.New("Invalid request", ErrorRequestInvalid, "request", err)
	}

	req.Header = r.Headers
	req.Header.Set(HeaderGoblDate, strconv.Itoa(int(time.Now().UTC().Unix())))

	if r.Client == nil {
		r.Client = &http.Client{CheckRedirect: checkRedirect}
	} else {
		r.Client.CheckRedirect = checkRedirect
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, goblerr.New("Unable to send request", ErrorRequestFailed, "request", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, goblerr.New("Unable to read response", ErrorRequestFailed, "request", err)
	}

	resp.Body.Close()

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, goblerr.New("Unable to unmarshal", ErrorRequestFailed, "request", err)
	}
	response.HTTPCode = resp.StatusCode

	return &response, nil
}

// do not allow redirects
func checkRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("Redirects not supported")
}
