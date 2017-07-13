package job

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/modification"
	"github.com/stretchr/testify/assert"
)

func TestRestore(t *testing.T) {
	assert := assert.New(t)

	r := &restore{
		job: model.Job{
			ID:   uuid.New().String(),
			Meta: &model.JobMeta{},
		},
		stateM:     &sync.Mutex{},
		MaxWorkers: 1,
	}

	f, err := createTestRestoreFile()
	if !assert.Nil(err) {
		return
	}

	r.coordClient = &tgrpc{rclient: &rclient{files: []*files.File{f}}}

	defer os.Remove("rtest.log")
	// local file save name will always be the same as the name is based on the fileSig
	defer os.Remove("435f25614eb9b7f51ab7921c5ff09992")

	r.job.Definition = &model.JobDefinition{
		Modifications: []modification.Definition{
			modification.Definition{Name: "compress", Options: map[string]string{"level": "5"}}},
		To: []engine.Definition{
			engine.Definition{
				Name:    engine.NameLogger,
				Options: map[string]string{engine.LoggerOptionLogPath: "rtest.log", engine.LoggerOptionOverwrite: "true"}}},
		From: &engine.Definition{
			Name:    engine.NameLocalFile,
			Options: map[string]string{engine.LocalFileOptionSavePath: "./", engine.LocalFileOptionOverwrite: "false"}},
	}

	finish := make(chan string)

	go r.Run(finish)

	id := <-finish
	assert.Equal(id, r.job.ID)

	fdata, err := ioutil.ReadFile("rtest.log")
	if !assert.Nil(err) {
		return
	}

	lines := strings.Split(string(fdata), "\n")
	if !assert.Len(lines, 2) {
		return
	}
}

func createTestRestoreFile() (*files.File, error) {
	// read data and compress
	data, err := ioutil.ReadFile("restore.go")
	if err != nil {
		return nil, err
	}

	var bbuf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&bbuf, 5)
	_, err = gz.Write(data)
	if err != nil {
		return nil, err
	}

	err = gz.Flush()
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	le := &engine.LocalFile{}
	le.ConfigureSave(map[string]string{engine.LocalFileOptionSavePath: "./", engine.LocalFileOptionOverwrite: "false"})
	errc := make(chan error, 1)

	file := files.File{
		Signature: files.Signature{Path: "test1", Hash: "asdf", Modifications: []string{"compress"}},
	}
	le.Save(bytes.NewReader(bbuf.Bytes()), file, errc)

	select {
	case err := <-errc:
		if err != nil {
			return nil, err
		}
	default:
		//continue
	}

	return &file, nil
}
