package leveldb

import (
	"testing"

	"github.com/sethjback/gobl/model"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	assert := assert.New(t)

	s, err := testDB()
	if !assert.Nil(err) {
		return
	}
	defer s.Close()

	err = s.SaveKey("test", model.Key{Key: "asdf", Certificate: "asdf"})
	assert.Nil(err)

	key, err := s.GetKey("test")
	assert.Nil(err)
	assert.Equal(&model.Key{Key: "asdf", Certificate: "asdf"}, key)
}
