package httpapi

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getURLVals() url.Values {
	v := url.Values{}
	v.Add("z", "21")
	v.Add("z", "12")
	v.Add("a", "1")
	v.Add("c", "1")
	v.Add("c", "az")
	return v
}

func TestQueryString(t *testing.T) {
	v := getURLVals()

	s := queryString(v)
	assert.Equal(t, "a=1&c=1&c=az&z=12&z=21", s, "Query string not sorted alpha")
}

func TestBodyHash(t *testing.T) {
	body := bytes.NewReader([]byte(`{"test1":"val1","test2":2}`))
	expected := "5f6757721a51d15bcf9d2efa81fd80bfd7ade0908e72f71f3324cf9bc5cc80c5"

	assert.Equal(t, expected, bodyHash(body))
}

func TestRequestString(t *testing.T) {
	h := http.Header{}
	h.Set(HeaderGoblDate, "1234")
	h.Set("authorization", "Bearer asdf.asdf.asdf")

	r := &Request{
		Headers: h,
		Body:    bytes.NewReader([]byte(`{"test1":"val1","test2":2}`)),
		Host:    "test.com",
		Path:    "/get/this",
		Method:  "POST",
		Query:   getURLVals(),
	}

	correct := "POST\ntest.com\n/get/this\na=1&c=1&c=az&z=12&z=21\nauthorization:Bearer asdf.asdf.asdf\nx-gobl-date:1234\n5f6757721a51d15bcf9d2efa81fd80bfd7ade0908e72f71f3324cf9bc5cc80c5"
	assert.Equal(t, correct, r.String())
}
