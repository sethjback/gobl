package engine

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sethjback/gobl/files"
	"github.com/stretchr/testify/assert"
)

func TestLocalFile(t *testing.T) {
	assert := assert.New(t)
	l := LocalFile{}

	err := l.ConfigureSave(map[string]string{"asdf": ""})
	assert.NotNil(err)

	err = l.ConfigureSave(map[string]string{"savePath": "true"})
	assert.NotNil(err)

	err = l.ConfigureSave(map[string]string{"overwrite": "123"})
	assert.NotNil(err)

	err = l.ConfigureSave(map[string]string{"savePath": "./", "overwrite": "false"})
	assert.Nil(err)

	file := files.File{Signature: files.Signature{Path: "/the/test/path/test1"}}

	fHash, err := hashFileSig(file.Signature)
	assert.Nil(err)

	save, err := l.ShouldSave(file)
	assert.True(save)
	assert.Nil(err)

	dataToSave := []byte("This is some test data to save that really should be saved")
	// buffered so the save routine doesn't block if there is an error
	errc := make(chan error, 3)

	l.Save(bytes.NewReader(dataToSave), file, errc)

	select {
	case e := <-errc:
		assert.Nil(e)
	default:
		//good
	}

	fData, err := ioutil.ReadFile(fHash)
	assert.Nil(err)
	assert.Equal(dataToSave, fData)

	save, err = l.ShouldSave(file)
	assert.False(save)
	assert.Nil(err)

	file.Modifications = append(file.Modifications, "test1")

	fHash2, err := hashFileSig(file.Signature)
	assert.Nil(err)

	save, err = l.ShouldSave(file)
	assert.True(save)
	assert.Nil(err)

	l.Save(bytes.NewReader(dataToSave), file, errc)

	select {
	case e := <-errc:
		assert.Nil(e)
	default:
		//good
	}

	fData, err = ioutil.ReadFile(fHash2)
	assert.Nil(err)
	assert.Equal(dataToSave, fData)

	err = os.Remove(fHash)
	assert.Nil(err)

	err = os.Remove(fHash2)
	assert.Nil(err)
}
