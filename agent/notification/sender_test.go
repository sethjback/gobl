package notification

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/keys"
	"github.com/stretchr/testify/assert"
)

type testNote struct {
	d string
	b []byte
}

func (ts testNote) Destination() string {
	return ts.d
}

func (ts testNote) Body() []byte {
	return ts.b
}

func TestSender(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bdy, err := ioutil.ReadAll(r.Body)
		assert.Nil(err)
		r.Body.Close()
		assert.NotNil(r.Header.Get(httpapi.HeaderGoblSig))
		if assert.NotEmpty(bdy) {
			switch string(bdy) {
			case "succeed":
				w.WriteHeader(200)
				fmt.Fprintln(w, `{"message":"success"}`)
			default:
				w.WriteHeader(400)
				fmt.Fprintln(w, `{"message":"invalid request"}`)
			}
		}

	}))
	defer ts.Close()

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if !assert.Nil(err) {
		return
	}

	s := &Sender{
		client:  &http.Client{},
		message: &Message{retry: 0, note: testNote{d: ts.URL, b: []byte("succeed")}},
		signer:  keys.NewSigner(pk),
	}

	r := s.Do().(*Result)

	assert.Equal(Success, r.state)

	s.message = &Message{retry: 0, note: testNote{d: ts.URL, b: []byte("fail")}}

	r = s.Do().(*Result)
	assert.Equal(Retry, r.state)
	assert.Equal(r.message, s.message)
}
