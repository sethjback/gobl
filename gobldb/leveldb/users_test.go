package leveldb

import (
	"testing"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/util/log"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	assert := assert.New(t)
	log.Init(config.Log{Level: log.Level.Error})
	s, err := testDB()
	if !assert.Nil(err) {
		return
	}
	defer s.Close()

	u := model.User{Email: "test@testing.com", Password: "asdf123"}
	u2 := model.User{Email: "test2@testing.com", Password: "asdf123"}
	u3 := model.User{Email: "test3@testing.com", Password: "asdf123"}

	err = s.SaveUser(u)
	assert.Nil(err)

	u1, err := s.GetUser(u.Email)
	if !assert.Nil(err) {
		return
	}
	assert.Equal(u, *u1)

	u.Password = "321fdsa"

	err = s.SaveUser(u)
	assert.Nil(err)

	u1, err = s.GetUser(u.Email)
	assert.Nil(err)
	assert.Equal(u, *u1)

	err = s.SaveUser(u2)
	assert.Nil(err)

	err = s.SaveUser(u3)
	assert.Nil(err)

	ulist, err := s.UserList()
	if assert.Nil(err) && assert.Len(ulist, 3) {
		assert.Contains(ulist, u)
		assert.Contains(ulist, u2)
		assert.Contains(ulist, u3)
	}
}
