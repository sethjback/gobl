package leveldb

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/modification"
	"github.com/sethjback/gobl/util/log"
	"github.com/stretchr/testify/assert"
)

func TestJobs(t *testing.T) {
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
	if !assert.Nil(err) {
		return
	}

	jd := model.JobDefinition{
		ID: uuid.New().String(),
		To: []engine.Definition{
			engine.Definition{Name: "test", Options: map[string]interface{}{"test": float64(1)}},
		},
		Modifications: []modification.Definition{
			modification.Definition{Name: "test", Options: map[string]interface{}{"test": float64(1)}},
		},
		Paths: []model.Path{
			model.Path{Root: "/dir1/dir2", Excludes: []string{"*.jpg"}},
		},
	}

	err = s.SaveJobDefinition(jd)
	if !assert.Nil(err) {
		return
	}

	j := model.Job{
		ID:         uuid.New().String(),
		Agent:      &a,
		Definition: &jd,
		Meta: &model.JobMeta{
			State: "running",
			Start: time.Now(),
			End:   time.Now(),
		},
	}

	err = s.SaveJob(j)
	assert.Nil(err)

	j1, err := s.GetJob(j.ID)
	if assert.Nil(err) {
		assert.Equal(j, *j1)
	}

	j.Meta.State = "finished"

	err = s.SaveJob(j)
	assert.Nil(err)

	j1, err = s.GetJob(j.ID)
	if assert.Nil(err) {
		assert.Equal(j, *j1)
	}

	a1 := model.Agent{
		ID:        uuid.New().String(),
		Name:      "Test Agent 1",
		Address:   "127.0.0.1:8080",
		PublicKey: "asdfasdfasdf",
	}

	err = s.SaveAgent(a1)
	if !assert.Nil(err) {
		return
	}

	err = s.SaveJob(model.Job{
		ID:         uuid.New().String(),
		Agent:      &a1,
		Definition: &jd,
		Meta: &model.JobMeta{
			State: "finished",
			Start: time.Date(2017, time.January, 1, 12, 12, 12, 0, time.UTC),
			End:   time.Date(2017, time.January, 2, 12, 12, 12, 0, time.UTC),
		},
	})
	assert.Nil(err)

	err = s.SaveJob(model.Job{
		ID:         uuid.New().String(),
		Agent:      &a1,
		Definition: &jd,
		Meta: &model.JobMeta{
			State: "finished",
			Start: time.Date(2017, time.January, 1, 12, 12, 12, 0, time.UTC),
			End:   time.Date(2017, time.January, 2, 12, 12, 12, 0, time.UTC),
		},
	})
	assert.Nil(err)

	err = s.SaveJob(model.Job{
		ID:         uuid.New().String(),
		Agent:      &a,
		Definition: &jd,
		Meta: &model.JobMeta{
			State: "finished",
			Start: time.Date(2017, time.January, 1, 12, 12, 12, 0, time.UTC),
			End:   time.Date(2017, time.January, 2, 12, 12, 12, 0, time.UTC),
		},
	})
	assert.Nil(err)

	err = s.SaveJob(model.Job{
		ID:         uuid.New().String(),
		Agent:      &a,
		Definition: &jd,
		Meta: &model.JobMeta{
			State: "finished",
			Start: time.Date(2017, time.February, 1, 12, 12, 12, 0, time.UTC),
			End:   time.Date(2017, time.February, 2, 12, 12, 12, 0, time.UTC),
		},
	})
	assert.Nil(err)

	err = s.SaveJob(model.Job{
		ID:         uuid.New().String(),
		Agent:      &a,
		Definition: &jd,
		Meta: &model.JobMeta{
			State: "running",
			Start: time.Now(),
			End:   time.Now(),
		},
	})
	assert.Nil(err)

	err = s.SaveJob(model.Job{
		ID:         uuid.New().String(),
		Agent:      &a1,
		Definition: &jd,
		Meta: &model.JobMeta{
			State: "running",
			Start: time.Now(),
			End:   time.Now(),
		},
	})
	assert.Nil(err)

	jds, err := s.JobList(map[string]string{"state": "finished"})
	if assert.Nil(err) {
		assert.Len(jds, 5)
	}

	jds, err = s.JobList(map[string]string{"state": "running"})
	if assert.Nil(err) {
		assert.Len(jds, 2)
	}

	jds, err = s.JobList(map[string]string{"agent": a.ID})
	if assert.Nil(err) {
		assert.Len(jds, 4)
	}

	jds, err = s.JobList(map[string]string{})
	if assert.Nil(err) {
		assert.Len(jds, 7)
	}

	jds, err = s.JobList(map[string]string{"limit": "2"})
	if assert.Nil(err) {
		assert.Len(jds, 2)
	}

	jds, err = s.JobList(map[string]string{"start": "2017-02-01 00:00"})
	if assert.Nil(err) {
		assert.Len(jds, 4)
	}

	jds, err = s.JobList(map[string]string{"start": "2017-02-01 00:00", "end": "2017-03-01 00:00"})
	if assert.Nil(err) {
		assert.Len(jds, 1)
	}

	jds, err = s.JobList(map[string]string{"end": "2017-02-01 00:00"})
	if assert.Nil(err) {
		assert.Len(jds, 3)
	}
}
