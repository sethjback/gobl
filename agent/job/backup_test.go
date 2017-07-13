package job

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/modification"
	"github.com/stretchr/testify/assert"
)

func TestBackup(t *testing.T) {
	assert := assert.New(t)

	b := &backup{
		job: model.Job{
			ID:   uuid.New().String(),
			Meta: &model.JobMeta{},
		},
		stateM:      &sync.Mutex{},
		MaxWorkers:  1,
		coordClient: &tgrpc{},
	}

	b.job.Definition = &model.JobDefinition{
		Paths: []model.Path{model.Path{Root: "test"}},
		Modifications: []modification.Definition{
			modification.Definition{Name: "compress", Options: map[string]string{"level": "5"}}},
		To: []engine.Definition{
			engine.Definition{
				Name:    engine.NameLogger,
				Options: map[string]string{engine.LoggerOptionLogPath: "btest.log", engine.LoggerOptionOverwrite: "false"}},
			engine.Definition{
				Name:    engine.NameLocalFile,
				Options: map[string]string{engine.LocalFileOptionSavePath: "saveDir", engine.LocalFileOptionOverwrite: "false"}}},
	}

	defer os.RemoveAll("btest.log")
	defer os.RemoveAll("saveDir")

	finish := make(chan string)

	go b.Run(finish)

	id := <-finish
	assert.Equal(id, b.job.ID)

	fdata, err := ioutil.ReadFile("btest.log")
	if !assert.Nil(err) {
		return
	}

	lines := strings.Split(string(fdata), "\n")
	if !assert.Len(lines, 21) {
		return
	}

	fcount := 0
	err = filepath.Walk("saveDir", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fcount++
		}
		return nil
	})

	assert.Nil(err)
	assert.Equal(20, fcount)
}

func TestBuildBackupFileList(t *testing.T) {
	assert := assert.New(t)

	c := make(chan struct{})
	errc := make(chan error, 20)
	path := model.Path{Root: "test"}
	in := buildBackupFileList(c, []model.Path{path}, errc)
	fCount := 0
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range in {
			fCount++
		}
	}()

	wg.Wait()
	close(errc)
	ecount := 0
	for range errc {
		ecount++
	}
	assert.Equal(0, ecount)
	assert.Equal(20, fCount)

	fCount = 0
	errc = make(chan error, 20)
	in = buildBackupFileList(c, []model.Path{path}, errc)
	wg = &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range in {
			fCount++
			if fCount == 10 {
				close(c)
			}
		}
	}()

	wg.Wait()
	assert.InDelta(10, fCount, 1)
}

func buildDirectoryTree() error {
	e := os.Mkdir("test", os.ModePerm)
	if e != nil {
		return e
	}

	for i := 0; i < 10; i++ {
		e = ioutil.WriteFile("test/tfile"+strconv.Itoa(i), []byte("Test file"), os.ModePerm)
		if e != nil {
			return e
		}

	}

	e = os.Mkdir("test/test1", os.ModePerm)
	if e != nil {
		return e
	}
	for i := 0; i < 10; i++ {
		e = ioutil.WriteFile("test/test1/tfile"+strconv.Itoa(i), []byte("Test file"), os.ModePerm)
		if e != nil {
			return e
		}

	}

	return nil
}

func cleanUpDirectoryTree() {
	os.RemoveAll("test")
}
