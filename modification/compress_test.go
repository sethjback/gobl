package modification

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/sethjback/gobl/goblerr"
	"github.com/stretchr/testify/assert"
)

func TestCompressConfigure(t *testing.T) {
	assert := assert.New(t)

	c := &Compress{}

	err := c.Configure(map[string]string{})
	assert.Nil(err)
	assert.Equal(5, c.level)
	assert.Equal("gzip", c.method)

	err = c.Configure(map[string]string{"method": "zlib"})
	if assert.NotNil(err) {
		gerr, ok := err.(goblerr.Error)
		if assert.True(ok) {
			assert.Equal(ErrorInvalidOptionValue, gerr.Code)
		}
	}

	err = c.Configure(map[string]string{"level": "10"})
	if assert.NotNil(err) {
		gerr, ok := err.(goblerr.Error)
		if assert.True(ok) {
			assert.Equal(ErrorInvalidOptionValue, gerr.Code)
		}
	}

	err = c.Configure(map[string]string{"level": "23"})
	if assert.NotNil(err) {
		gerr, ok := err.(goblerr.Error)
		if assert.True(ok) {
			assert.Equal(ErrorInvalidOptionValue, gerr.Code)
		}
	}

	err = c.Configure(map[string]string{"level": "2", "method": "gzip"})
	assert.Nil(err)
	assert.Equal(2, c.level)
	assert.Equal("gzip", c.method)
}

func TestCompress(t *testing.T) {
	assert := assert.New(t)

	c := &Compress{}
	c.Configure(map[string]string{})
	c.Direction(Forward)

	toCompress := []byte("this is the string to test compress")
	errc := make(chan error, 2)

	compressed, err := ioutil.ReadAll(c.Process(bytes.NewReader(toCompress), errc))
	assert.Nil(err)
	assert.NotEmpty(compressed)

	c.Direction(Backward)
	restored, err := ioutil.ReadAll(c.Process(bytes.NewReader(compressed), errc))

	assert.Nil(err)
	assert.Equal(restored, toCompress)
}
