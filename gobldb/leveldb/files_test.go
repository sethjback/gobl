package leveldb

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/util/log"
	"github.com/stretchr/testify/assert"
)

func TestFiles(t *testing.T) {
	assert := assert.New(t)
	log.Init(config.Log{Level: log.Level.Error})

	s, err := testDB()
	if !assert.Nil(err) {
		return
	}
	defer s.Close()
	jobID := uuid.New().String()

	f := model.JobFile{
		State: "complete",
		File: files.File{
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
	}

	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}

	f1, err := s.getFile(jobID, "/dir1/dir2/dir3/file1.jpg")
	if assert.Nil(err) {
		assert.Equal(f.File, f1.File)
		assert.Equal(f.State, f1.State)
		assert.Equal(f.Error, f1.Error)
	}

	f.State = "failed"
	f.Error = "unable to write file: device busy"
	err = s.SaveJobFile(jobID, f)
	assert.Nil(err)

	f1, err = s.getFile(jobID, "/dir1/dir2/dir3/file1.jpg")
	if assert.Nil(err) {
		assert.Equal(f.File, f1.File)
		assert.Equal(f.State, f1.State)
		assert.Equal(f.Error, f1.Error)
	}

	f.File.Path = "/dir1/dir2/dir3/file2.jpg"
	f.State = "success"
	f.Error = ""
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}
	f1, err = s.getFile(jobID, "/dir1/dir2/dir3/file2.jpg")
	if assert.Nil(err) {
		assert.Equal(f.File, f1.File)
		assert.Equal(f.State, f1.State)
		assert.Equal(f.Error, f1.Error)
	}

	f.File.Path = "/dir1/dir2/dir3/file3.jpg"
	f.State = "success"
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}
	f1, err = s.getFile(jobID, "/dir1/dir2/dir3/file3.jpg")
	if assert.Nil(err) {
		assert.Equal(f.File, f1.File)
		assert.Equal(f.State, f1.State)
		assert.Equal(f.Error, f1.Error)
	}

	f.File.Path = "/dir1/dir2/dir3.2/file1.jpg"
	f.State = "success"
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}
	f.File.Path = "/dir1/dir2/dir3.2/file2.jpg"
	f.State = "success"
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}

	f.File.Path = "/dir1/dir2/dir3.3/file1.jpg"
	f.State = "failed"
	f.Error = "unable to save file (device busy)"
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}
	f.File.Path = "/dir1/dir2/dir3.3/file2.jpg"
	f.State = "success"
	f.Error = ""
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}
	f.File.Path = "/dir1/dir2/dir3.3/file3.jpg"
	f.State = "failed"
	f.Error = "unable to save file (device busy)"
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}

	jfs, err := s.JobFileList(jobID, map[string]string{"state": "failed"})
	if assert.Nil(err) {
		assert.Len(jfs, 3)
	}

	jfs, err = s.JobFileList(jobID, map[string]string{"state": "failed", "dir": "/dir1/dir2/dir3"})
	if assert.Nil(err) {
		assert.Len(jfs, 1)
	}

	jfs, err = s.JobFileList(jobID, map[string]string{"state": "failed", "dir": "/dir1/dir2/dir3.2"})
	if assert.Nil(err) {
		assert.Len(jfs, 0)
	}

	jfs, err = s.JobFileList(jobID, map[string]string{"dir": "/dir1/dir2/dir3"})
	if assert.Nil(err) {
		assert.Len(jfs, 3)
	}

	jfc, err := s.jobFileCount(jobID, map[string]string{"state": "failed"})
	if assert.Nil(err) {
		assert.Equal(3, jfc)
	}

	jfc, err = s.jobFileCount(jobID, map[string]string{"dir": "*"})
	if assert.Nil(err) {
		assert.Equal(8, jfc)
	}

	f.File.Path = "/dir1.2/dir2/dir3.2/file2.jpg"
	f.State = "success"
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}

	f.File.Path = "/dir1.2/dir2.2/dir3.2/file2.jpg"
	f.State = "success"
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}
	f.File.Path = "/dir1.2/dir2.3/dir3.1/file2.jpg"
	f.State = "success"
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}
	f.File.Path = "/dir1.2/dir2.3/dir3.2/file2.jpg"
	f.State = "success"
	err = s.SaveJobFile(jobID, f)
	if !assert.Nil(err) {
		return
	}

	dirs, err := s.JobDirectories(jobID, "/dir1.2")
	if assert.Nil(err) {
		assert.Len(dirs, 3)
		assert.Contains(dirs, "dir2")
		assert.Contains(dirs, "dir2.2")
		assert.Contains(dirs, "dir2.3")
	}

	dirs, err = s.JobDirectories(jobID, "/dir1.2/dir2")
	if assert.Nil(err) {
		assert.Len(dirs, 1)
		assert.Contains(dirs, "dir3.2")
	}

	dirs, err = s.JobDirectories(jobID, "/dir1.2/dir2.3")
	if assert.Nil(err) {
		assert.Len(dirs, 2)
		assert.Contains(dirs, "dir3.1")
		assert.Contains(dirs, "dir3.2")
	}

	dirs, err = s.JobDirectories(jobID, "/")
	if assert.Nil(err) {
		assert.Len(dirs, 2)
		assert.Contains(dirs, "dir1")
		assert.Contains(dirs, "dir1.2")
	}
}
