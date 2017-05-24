package leveldb

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/util/log"
	"github.com/stretchr/testify/assert"
)

func TestSchedule(t *testing.T) {
	assert := assert.New(t)
	log.Init(config.Log{Level: log.Level.Error})
	s, err := testDB()
	if !assert.Nil(err) {
		return
	}

	defer s.Close()

	sc := model.Schedule{
		ID:              uuid.New().String(),
		JobDefinitionID: uuid.New().String(),
		AgentID:         uuid.New().String(),
		Seconds:         "0",
		Minutes:         "0",
		Hour:            "0",
		DOM:             "*",
		MON:             "*",
		DOW:             "*",
	}

	err = s.SaveSchedule(sc)
	assert.Nil(err)

	sc1, err := s.GetSchedule(sc.ID)
	if assert.Nil(err) {
		assert.Equal(sc, *sc1)
	}

	sc.Seconds = "3"
	sc.Minutes = "4"

	err = s.SaveSchedule(sc)
	assert.Nil(err)

	sc1, err = s.GetSchedule(sc.ID)
	if assert.Nil(err) {
		assert.Equal(sc, *sc1)
	}

	sc2 := model.Schedule{
		ID:              uuid.New().String(),
		JobDefinitionID: uuid.New().String(),
		AgentID:         uuid.New().String(),
		Seconds:         "0",
		Minutes:         "0",
		Hour:            "0",
		DOM:             "*",
		MON:             "*",
		DOW:             "*",
	}

	err = s.SaveSchedule(sc2)
	assert.Nil(err)

	sc3 := model.Schedule{
		ID:              uuid.New().String(),
		JobDefinitionID: uuid.New().String(),
		AgentID:         uuid.New().String(),
		Seconds:         "0",
		Minutes:         "0",
		Hour:            "0",
		DOM:             "*",
		MON:             "*",
		DOW:             "*",
	}

	err = s.SaveSchedule(sc3)
	assert.Nil(err)

	sch, err := s.ScheduleList()
	if assert.Nil(err) {
		assert.Len(sch, 3)
		assert.Contains(sch, sc)
		assert.Contains(sch, sc2)
		assert.Contains(sch, sc3)
	}

	err = s.DeleteSchedule(sc3.ID)
	assert.Nil(err)

	sch, err = s.ScheduleList()
	if assert.Nil(err) {
		assert.Len(sch, 2)
		assert.Contains(sch, sc)
		assert.Contains(sch, sc2)
		assert.NotContains(sch, sc3)
	}
}
