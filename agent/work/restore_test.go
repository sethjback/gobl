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
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/modification"
	"github.com/stretchr/testify/assert"
)

func TestRestore(t *testing.T) {
	assert := assert.New(t)

	// read data and compress
	data, err := ioutil.ReadFile("restore.go")
	if !assert.Nil(err) {
		return
	}

	var bbuf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&bbuf, 5)
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

	le := &engine.LocalFile{}
	le.ConfigureSave(map[string]interface{}{engine.LocalFileOptionSavePath: "./", engine.LocalFileOptionOverwrite: false})
	errc := make(chan error, 1)

	file := files.File{
		Signature: files.Signature{Path: "test1", Hash: "asdf", Modifications: []string{"compress"}},
	}
	le.Save(bytes.NewReader(bbuf.Bytes()), file, errc)

	select {
	case err := <-errc:
		if !assert.Nil(err) {
			return
		}
	default:
		//continue
	}

	close(errc)

	r := Restore{
		File:          file,
		Modifications: []modification.Definition{modification.Definition{Name: "compress"}},
		From:          engine.Definition{Name: "localfile", Options: map[string]interface{}{engine.LocalFileOptionSavePath: "./", engine.LocalFileOptionOverwrite: false}},
		To: []engine.Definition{
			engine.Definition{
				Name:    engine.NameLogger,
				Options: map[string]interface{}{engine.LoggerOptionLogPath: "rtest.log", engine.LoggerOptionOverwrite: true}}},
	}

	res := r.Do()
	if assert.NotNil(res) {
		jf, ok := res.(model.JobFile)
		if assert.True(ok) && assert.Equal(StateComplete, jf.State) {
			// Hack: allow the file to flush to disk
			<-time.NewTimer(1 * time.Second).C
			fdata, err := ioutil.ReadFile("rtest.log")
			if !assert.Nil(err) {
				os.Remove("rtest.log")
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

	os.Remove("rtest.log")
	// local file save name will always be the same as the name is based on the fileSig
	os.Remove("435f25614eb9b7f51ab7921c5ff09992")
}
