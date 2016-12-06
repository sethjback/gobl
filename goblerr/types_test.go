package goblerr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseError(t *testing.T) {
	assert := assert.New(t)

	be := newBaseError("Test Message", "TESTCODE", nil, map[string]string{"test": "ing"})
	assert.NotNil(be)
	assert.Equal("Test Message", be.Message())
	assert.Equal("TESTCODE", be.Code())
	assert.Equal(map[string]string{"test": "ing"}, be.Detail())

	b2 := newBaseError("Test 2", "t2", be, nil)

	assert.Equal(be, b2.Origin())
}
