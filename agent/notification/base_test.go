package notification

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/util/log"
	"github.com/stretchr/testify/assert"
)

type testNotification struct {
	d string
	b []byte
}

func (ts testNotification) Destination() string {
	return ts.d
}

func (ts testNotification) Body() []byte {
	return ts.b
}

func TestNotificationQueue(t *testing.T) {
	assert := assert.New(t)

	log.Init(config.Log{Level: log.Level.Warn})

	msgCount := int64(0)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bdy, err := ioutil.ReadAll(r.Body)
		assert.Nil(err)
		r.Body.Close()
		assert.NotNil(r.Header.Get(httpapi.HeaderGoblSig))
		if assert.NotEmpty(bdy) {
			switch string(bdy) {
			case "succeed":
				atomic.AddInt64(&msgCount, 1)
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

	n := newBase(nil, pk)
	n.Start()
	for i := 0; i < 10; i++ {
		n.Send(testNotification{d: ts.URL, b: []byte("fail")})
	}

	n.Stop()
	assert.Equal(10, n.pending.Length()+n.retry.Length())

	msgCount = 0

	n = newBase(nil, pk)
	n.Start()
	for i := 0; i < 10; i++ {
		n.Send(testNotification{d: ts.URL, b: []byte("succeed")})
	}

	n.Stop()
	assert.Equal(10, n.pending.Length()+int(msgCount))

	msgCount = 0

	n = newBase(nil, pk)
	n.Start()

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		for i := 0; i < 10; i++ {
			n.Send(testNotification{d: ts.URL, b: []byte("succeed")})
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10; i++ {
			n.Send(testNotification{d: ts.URL, b: []byte("fail")})
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10; i++ {
			n.Send(testNotification{d: ts.URL, b: []byte("succeed")})
		}
		wg.Done()
	}()

	wg.Wait()

	n.Stop()
	assert.Equal(30, n.retry.Length()+n.pending.Length()+int(msgCount))
}
