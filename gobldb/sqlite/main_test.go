package sqlite

import (
	"testing"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/util/log"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	assert := assert.New(t)
	log.Init(config.Log{Level: log.Level.Error})

	sql := SQLite{}
	if assert.Nil(sql.Init(config.DB{Path: ""})) {
		sql.Close()
	}
}

func testDB() (*SQLite, error) {
	s := &SQLite{}
	err := s.Init(config.DB{Path: ""})
	return s, err
}
