package gobldb

import (
	"testing"

	"github.com/sethjback/gobl/config"
	"github.com/stretchr/testify/assert"
)

func TestSaveConfig(t *testing.T) {
	assert := assert.New(t)

	cs := config.New()

	assert.Nil(SaveConfig(cs, map[string]string{"DB_PATH": "", "DB_DRIVER": "invalid"}))

	dbc := configFromStore(cs)
	if assert.NotNil(dbc) {
		assert.Equal("invalid", dbc.driver)
	}

	db, err := Get(cs)
	assert.Nil(db)
	assert.NotNil(err)

	cs = config.New()

	assert.Nil(SaveConfig(cs, map[string]string{"DB_PATH": "", "DB_DRIVER": "leveldb"}))

	dbc = configFromStore(cs)
	if assert.NotNil(dbc) {
		assert.Equal("leveldb", dbc.driver)
	}

	db, err = Get(cs)
	assert.Nil(err)
	if assert.NotNil(db) {
		db.Close()
	}
}
