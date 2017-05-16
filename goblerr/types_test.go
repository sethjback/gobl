package goblerr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseError(t *testing.T) {
	assert := assert.New(t)

	be := newBaseError("Test Message", "TESTCODE", map[string]string{"test": "ing"})
	assert.NotNil(be)
	assert.Equal("Test Message", be.Message())
	assert.Equal("TESTCODE", be.Code())
	assert.Equal(map[string]string{"test": "ing"}, be.Detail())
}
