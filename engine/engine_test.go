package engine

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/sethjback/gobl/files"
	"github.com/stretchr/testify/assert"
)

type TestEngine struct {
	sMutex *sync.Mutex
	saved  []byte
}

func (t *TestEngine) GetSaved() []byte {
	t.sMutex.Lock()
	s := t.saved
	t.sMutex.Unlock()
	return s
}

func (t *TestEngine) Save(input io.Reader, file files.File, errc chan<- error) {
	if file.Path == "failSave" {
		errc <- errors.New("Failed")
		return
	}

	t.sMutex.Lock()
	b, e := ioutil.ReadAll(input)
	t.saved = b
	t.sMutex.Unlock()
	if e != nil {
		errc <- errors.New("Failed")
		return
	}
}
func (t *TestEngine) Retrieve(file files.File) (io.Reader, error) {
	return nil, nil
}
func (t *TestEngine) ShouldSave(file files.File) (bool, error) {
	if file.Path == "failShouldSave" {
		return false, errors.New("Fail")
	}
	return true, nil
}
func (t *TestEngine) Name() string {
	return "TestEngine"
}
func (t *TestEngine) SaveOptions() []Option {
	return nil
}
func (t *TestEngine) ConfigureSave(map[string]interface{}) error {
	return nil
}

func (t *TestEngine) Restore(input io.Reader, file files.File, errc chan<- error) {
	if file.Path == "failRestore" {
		errc <- errors.New("Failed")
		return
	}

	t.sMutex.Lock()
	b, e := ioutil.ReadAll(input)
	t.saved = b
	t.sMutex.Unlock()
	if e != nil {
		errc <- errors.New("Failed")
		return
	}
}
func (t *TestEngine) ShouldRestore(file files.File) (bool, error) {
	if file.Path == "failShouldRestore" {
		return false, errors.New("Fail")
	}
	return true, nil
}

func (t *TestEngine) RestoreOptions() []Option {
	return nil
}
func (t *TestEngine) ConfigureRestore(map[string]interface{}) error {
	return nil
}

func TestBackupEngine(t *testing.T) {
	assert := assert.New(t)

	t1 := &TestEngine{sMutex: &sync.Mutex{}}

	file := files.File{Signature: files.Signature{Path: "failShouldSave"}}

	egn, ok, err := NewBackupEngine(file, t1)
	assert.Nil(egn)
	assert.False(ok)
	assert.NotNil(err)

	file.Path = "failSave"

	toSave := []byte("This is the data that we want to save")
	donec := make(chan bool)

	egn, ok, err = NewBackupEngine(file, t1)
	assert.True(ok)
	go func() {
		_, cerr := io.Copy(egn, bytes.NewReader(toSave))
		donec <- assert.Nil(cerr)
	}()
	errc := egn.ErrorChan()

	select {
	case e := <-errc:
		assert.NotNil(e)
	case d := <-donec:
		assert.False(d)
		egn.Finish()
	}

	t2 := &TestEngine{sMutex: &sync.Mutex{}}
	t3 := &TestEngine{sMutex: &sync.Mutex{}}
	t4 := &TestEngine{sMutex: &sync.Mutex{}}

	file.Path = "asdf"

	egn, ok, err = NewBackupEngine(file, t1, t2, t3, t4)
	assert.True(ok)
	go func() {
		_, cerr := io.Copy(egn, bytes.NewReader(toSave))
		donec <- assert.Nil(cerr)
	}()

	errc = egn.ErrorChan()

	select {
	case e := <-errc:
		assert.Nil(e)
	case d := <-donec:
		assert.True(d)
	}

	egn.Finish()

	assert.Equal(toSave, t1.GetSaved())
	assert.Equal(toSave, t2.GetSaved())
	assert.Equal(toSave, t3.GetSaved())
	assert.Equal(toSave, t4.GetSaved())

}

func TestRestoreEngine(t *testing.T) {
	assert := assert.New(t)

	t1 := &TestEngine{sMutex: &sync.Mutex{}}

	file := files.File{Signature: files.Signature{Path: "failShouldRestore"}}

	egn, err := NewRestoreEngine(file, t1)
	assert.Nil(egn)
	assert.NotNil(err)

	file.Path = "failRestore"

	toSave := []byte("This is the data that we want to restore")
	donec := make(chan bool)

	egn, err = NewRestoreEngine(file, t1)
	go func() {
		_, cerr := io.Copy(egn, bytes.NewReader(toSave))
		donec <- assert.Nil(cerr)
	}()
	errc := egn.ErrorChan()

	select {
	case e := <-errc:
		assert.NotNil(e)
	case d := <-donec:
		assert.False(d)
		egn.Finish()
	}

	t2 := &TestEngine{sMutex: &sync.Mutex{}}
	t3 := &TestEngine{sMutex: &sync.Mutex{}}
	t4 := &TestEngine{sMutex: &sync.Mutex{}}

	file.Path = "asdf"

	egn, err = NewRestoreEngine(file, t1, t2, t3, t4)
	go func() {
		_, cerr := io.Copy(egn, bytes.NewReader(toSave))
		donec <- assert.Nil(cerr)
	}()

	errc = egn.ErrorChan()

	select {
	case e := <-errc:
		assert.Nil(e)
	case d := <-donec:
		assert.True(d)
	}

	egn.Finish()

	assert.Equal(toSave, t1.GetSaved())
	assert.Equal(toSave, t2.GetSaved())
	assert.Equal(toSave, t3.GetSaved())
	assert.Equal(toSave, t4.GetSaved())

}
