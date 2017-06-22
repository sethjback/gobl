package leveldb

import (
	"os"
	"testing"

	"github.com/sethjback/gobl/config"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)

	l, err := New("")

	assert.Nil(err)
	assert.Nil(l.Close())

	l, err = New("./test.ldb")

	assert.Nil(err)
	assert.Nil(l.Close())

	_, err = os.Stat("./test.ldb")
	assert.Nil(err)
	assert.Nil(os.RemoveAll("./test.ldb"))
}

func testDB() (*Leveldb, error) {
	return New("")
}

func TestSaveConfig(t *testing.T) {
	assert := assert.New(t)

	cs := config.New()
	l := &Leveldb{}

	assert.Nil(l.SaveConfig(cs, map[string]string{"DB_PATH": "./test"}))

	dbc := configFromStore(cs)
	if assert.NotNil(dbc) {
		assert.Equal("./test", dbc.path)
	}

	assert.Nil(l.SaveConfig(cs, map[string]string{}))

	dbc = configFromStore(cs)
	if assert.NotNil(dbc) {
		assert.Empty(dbc.path)
	}

}
