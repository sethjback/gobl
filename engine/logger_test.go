package engine

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sethjback/gobl/files"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	assert := assert.New(t)

	l := Logger{}

	err := l.ConfigureSave(map[string]string{"asdf": ""})
	assert.NotNil(err)

	err = l.ConfigureSave(map[string]string{"logPath": "true"})
	assert.NotNil(err)

	err = l.ConfigureSave(map[string]string{"overwrite": "123"})
	assert.NotNil(err)

	err = l.ConfigureSave(map[string]string{"logPath": "./testLog", "overwrite": "false"})
	assert.Nil(err)

	dataToSave := []byte("This is some test data to save that really should be saved")
	// buffered so the save routine doesn't block if there is an error
	errc := make(chan error, 3)

	file := files.File{Signature: files.Signature{Path: "/the/test/path/test1"}}

	l.Save(bytes.NewReader(dataToSave), file, errc)

	select {
	case e := <-errc:
		assert.Nil(e)
	default:
		//good
	}

	fData, err := ioutil.ReadFile("./testLog")
	assert.Nil(err)

	lline := &LogLine{}
	err = json.Unmarshal(fData, lline)
	assert.Nil(err)
	assert.Equal(file.Path, lline.File)
	assert.Equal(len(dataToSave), int(lline.Bytes))

	// make sure we don't overwrite
	err = l.ConfigureSave(map[string]string{"logPath": "./testLog", "overwrite": "false"})
	assert.Nil(err)

	dataToSave2 := []byte("This is some test data to save that really should be saved and should be longer")

	l.Save(bytes.NewReader(dataToSave2), file, errc)

	select {
	case e := <-errc:
		assert.Nil(e)
	default:
		//good
	}

	var ll []LogLine
	f, err := os.Open("./testLog")
	if assert.Nil(err) {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lline = &LogLine{}
			err = json.Unmarshal(scanner.Bytes(), &lline)
			if assert.Nil(err) {
				ll = append(ll, *lline)
			}
		}
		f.Close()
	}

	if assert.Len(ll, 2) {
		assert.Equal(len(dataToSave), int(ll[0].Bytes))
		assert.Equal(len(dataToSave2), int(ll[1].Bytes))
	}

	// make sure we don't overwrite
	err = l.ConfigureSave(map[string]string{"logPath": "./testLog", "overwrite": "true"})
	assert.Nil(err)

	l.Save(bytes.NewReader(dataToSave), file, errc)

	select {
	case e := <-errc:
		assert.Nil(e)
	default:
		//good
	}

	fData, err = ioutil.ReadFile("./testLog")
	assert.Nil(err)

	lline = &LogLine{}
	err = json.Unmarshal(fData, &lline)
	assert.Nil(err)
	assert.Equal(file.Path, lline.File)
	assert.Equal(len(dataToSave), int(lline.Bytes))

	err = os.Remove("./testLog")
	assert.Nil(err)
}
