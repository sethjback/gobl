package sqlite

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/modification"
	"github.com/sethjback/gobl/util/log"
	"github.com/stretchr/testify/assert"
)

func TestJobDefinitions(t *testing.T) {
	assert := assert.New(t)
	log.Init(config.Log{Level: log.Level.Error})

	s, err := testDB()
	if !assert.Nil(err) {
		return
	}
	defer s.Close()

	jd := model.JobDefinition{
		ID: uuid.New().String(),
		To: []engine.Definition{
			engine.Definition{Name: "test", Options: map[string]interface{}{"test": float64(1)}},
			engine.Definition{Name: "test1", Options: map[string]interface{}{"test2": "three"}},
		},
		From: &engine.Definition{Name: "test1", Options: map[string]interface{}{"test2": "three"}},
		Modifications: []modification.Definition{
			modification.Definition{Name: "test", Options: map[string]interface{}{"test": float64(1)}},
			modification.Definition{Name: "test1", Options: map[string]interface{}{"test2": "three"}},
		},
		Paths: []model.Path{
			model.Path{Root: "/dir1/dir2", Excludes: []string{"*.jpg"}},
			model.Path{Root: "/dir3/dir4", Excludes: []string{"*.png", "*.mov"}},
		},
		Files: []files.File{
			files.File{
				Meta: files.Meta{
					Mode: 111,
					UID:  1,
					GID:  1,
				},
				Signature: files.Signature{
					Path:          "/dir1/dir2/dir3/file1.jpg",
					Hash:          "asdf",
					Modifications: []string{"mod1", "mod2"},
				},
			},
		},
	}

	err = s.SaveJobDefinition(jd)
	if !assert.Nil(err) {
		return
	}

	jd1, err := s.GetJobDefinition(jd.ID)
	if assert.Nil(err) {
		assert.Equal(*jd1, jd)
	}

	jd.From = nil
	jd.To = nil
	jd.Paths = nil
	jd.Modifications = nil
	jd.Files = nil

	err = s.SaveJobDefinition(jd)
	if !assert.Nil(err) {
		return
	}

	jd1, err = s.GetJobDefinition(jd.ID)
	if assert.Nil(err) {
		assert.Equal(*jd1, jd)
	}

	err = s.DeleteJobDefinition(jd.ID)
	assert.Nil(err)

	jd1, err = s.GetJobDefinition(jd.ID)
	assert.Nil(jd1)
	if assert.NotNil(err) {
		assert.Equal("No job definition with that ID", err.Error())
	}

	err = s.SaveJobDefinition(jd)
	if !assert.Nil(err) {
		return
	}

	jd2 := model.JobDefinition{
		ID: uuid.New().String(),
		To: []engine.Definition{
			engine.Definition{Name: "test", Options: map[string]interface{}{"test": float64(1)}},
			engine.Definition{Name: "test1", Options: map[string]interface{}{"test2": "three"}},
		},
		From: &engine.Definition{Name: "test1", Options: map[string]interface{}{"test2": "three"}},
		Modifications: []modification.Definition{
			modification.Definition{Name: "test", Options: map[string]interface{}{"test": float64(1)}},
			modification.Definition{Name: "test1", Options: map[string]interface{}{"test2": "three"}},
		},
		Paths: []model.Path{
			model.Path{Root: "/dir1/dir2", Excludes: []string{"*.jpg"}},
			model.Path{Root: "/dir3/dir4", Excludes: []string{"*.png", "*.mov"}},
		},
		Files: []files.File{
			files.File{
				Meta: files.Meta{
					Mode: 111,
					UID:  1,
					GID:  1,
				},
				Signature: files.Signature{
					Path:          "/dir1/dir2/dir3/file1.jpg",
					Hash:          "asdf",
					Modifications: []string{"mod1", "mod2"},
				},
			},
		},
	}

	err = s.SaveJobDefinition(jd2)
	if !assert.Nil(err) {
		return
	}

	jd3 := model.JobDefinition{
		ID:   uuid.New().String(),
		From: &engine.Definition{Name: "test1", Options: map[string]interface{}{"test2": "three"}},
		Modifications: []modification.Definition{
			modification.Definition{Name: "test", Options: map[string]interface{}{"test": float64(1)}},
		},
		Paths: []model.Path{
			model.Path{Root: "/dir3/dir4", Excludes: []string{"*.png", "*.mov"}},
		},
	}

	err = s.SaveJobDefinition(jd3)
	if !assert.Nil(err) {
		return
	}

	jds, err := s.GetJobDefinitions()
	assert.Nil(err)
	if assert.Len(jds, 3) {
		assert.Contains(jds, jd)
		assert.Contains(jds, jd2)
		assert.Contains(jds, jd3)
	}

}
