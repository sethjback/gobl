package leveldb

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/util/log"
	"github.com/stretchr/testify/assert"
)

func TestAgents(t *testing.T) {
	assert := assert.New(t)
	log.Init(config.Log{Level: log.Level.Error})

	s, err := testDB()
	if !assert.Nil(err) {
		return
	}
	defer s.Close()
	a := model.Agent{
		ID:        uuid.New().String(),
		Name:      "Test Agent 1",
		Address:   "127.0.0.1:8080",
		PublicKey: "asdfasdfasdf",
	}

	err = s.SaveAgent(a)
	assert.Nil(err)

	a1, err := s.GetAgent(a.ID)
	assert.Nil(err)
	assert.Equal(*a1, a)

	a.Name = "Different"
	a.Address = "1.1.1.1:8123"
	a.PublicKey = ""

	err = s.SaveAgent(a)
	assert.Nil(err)

	a1, err = s.GetAgent(a.ID)
	assert.Nil(err)
	assert.Equal(*a1, a)

	a2 := model.Agent{
		ID:        uuid.New().String(),
		Name:      "Test Agent 2",
		Address:   "127.0.0.1:8080",
		PublicKey: "asdfasdfasdf",
	}
	a3 := model.Agent{
		ID:        uuid.New().String(),
		Name:      "Test Agent 3",
		Address:   "127.0.0.2:8080",
		PublicKey: "asdfasdfasdf",
	}
	a4 := model.Agent{
		ID:        uuid.New().String(),
		Name:      "Test Agent 4",
		Address:   "127.0.0.3:8080",
		PublicKey: "asdfasdfasdf",
	}

	err = s.SaveAgent(a2)
	assert.Nil(err)
	err = s.SaveAgent(a3)
	assert.Nil(err)
	err = s.SaveAgent(a4)
	assert.Nil(err)

	alist, err := s.AgentList()
	assert.Nil(err)
	assert.Len(alist, 4)

	assert.Contains(alist, a)
	assert.Contains(alist, a2)
	assert.Contains(alist, a3)
	assert.Contains(alist, a4)
}
