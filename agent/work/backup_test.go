package work

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/modification"
	"github.com/stretchr/testify/assert"
)

func TestBackup(t *testing.T) {
	assert := assert.New(t)

	// read data and compress
	data, err := ioutil.ReadFile("backup.go")
	if !assert.Nil(err) {
		return
	}
	var bbuf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&bbuf, 5)
	if !assert.Nil(err) {
		return
	}

	_, err = gz.Write(data)
	if !assert.Nil(err) {
		return
	}

	err = gz.Flush()
	if !assert.Nil(err) {
		return
	}

	err = gz.Close()
	if !assert.Nil(err) {
		return
	}

	data = bbuf.Bytes()

	b := Backup{
		File: "backup.go",
		Modifications: []modification.Definition{
			modification.Definition{
				Name:    "compress",
				Options: map[string]interface{}{"level": 5}}},
		Engines: []engine.Definition{
			engine.Definition{
				Name:    engine.NameLogger,
				Options: map[string]interface{}{engine.LoggerOptionLogPath: "btest.log", engine.LoggerOptionOverwrite: true}}},
	}

	r := b.Do()
	if assert.NotNil(r) {
		jf, ok := r.(model.JobFile)
		if assert.True(ok) && assert.Equal(StateComplete, jf.State) {
			// Hack: allow the file to flush to disk
			<-time.NewTimer(1 * time.Second).C
			fdata, err := ioutil.ReadFile("btest.log")
			if !assert.Nil(err) {
				os.Remove("btest.log")
				return
			}

			lines := strings.Split(string(fdata), "\n")
			if assert.Len(lines, 2) {
				var lline engine.LogLine
				err = json.Unmarshal([]byte(lines[0]), &lline)
				if assert.Nil(err) {
					assert.Equal(len(data), lline.Bytes)
				}
			}
		}
	}

	os.Remove("btest.log")
}
