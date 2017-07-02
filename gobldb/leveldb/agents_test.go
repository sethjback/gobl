package leveldb

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/model"
	"github.com/stretchr/testify/assert"
)

func TestAgents(t *testing.T) {
	assert := assert.New(t)

	s, err := testDB()
	if !assert.Nil(err) {
		return
	}
	defer s.Close()
	a := model.Agent{
		ID:      uuid.New().String(),
		Name:    "Test Agent 1",
		Address: "127.0.0.1:8080",
		Key:     &model.Key{Key: "asdf", Certificate: "asdf"},
	}

	err = s.SaveAgent(a)
	assert.Nil(err)

	a1, err := s.GetAgent(a.ID)
	assert.Nil(err)
	assert.Equal(*a1, a)

	a.Name = "Different"
	a.Address = "1.1.1.1:8123"
	a.Key = nil

	err = s.SaveAgent(a)
	assert.Nil(err)

	a1, err = s.GetAgent(a.ID)
	assert.Nil(err)
	assert.Equal(*a1, a)

	a2 := model.Agent{
		ID:      uuid.New().String(),
		Name:    "Test Agent 2",
		Address: "127.0.0.1:8080",
		Key:     &model.Key{Key: "asdf", Certificate: "asdf"},
	}
	a3 := model.Agent{
		ID:      uuid.New().String(),
		Name:    "Test Agent 3",
		Address: "127.0.0.2:8080",
		Key:     &model.Key{Key: "asdf", Certificate: "asdf"},
	}
	a4 := model.Agent{
		ID:      uuid.New().String(),
		Name:    "Test Agent 4",
		Address: "127.0.0.3:8080",
		Key:     &model.Key{Key: "asdf", Certificate: "asdf"},
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
