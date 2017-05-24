package leveldb

import (
	"os"
	"testing"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/util/log"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)
	log.Init(config.Log{Level: log.Level.Error})

	l, err := New(config.DB{Path: ""})

	assert.Nil(err)
	assert.Nil(l.Close())

	l, err = New(config.DB{Path: "./test.ldb"})

	assert.Nil(err)
	assert.Nil(l.Close())

	_, err = os.Stat("./test.ldb")
	assert.Nil(err)
	assert.Nil(os.RemoveAll("./test.ldb"))
}

func testDB() (*Leveldb, error) {
	return New(config.DB{Path: ""})
}
