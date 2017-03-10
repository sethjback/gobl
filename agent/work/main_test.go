package work

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"sync"

	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/files"
)

type TestEngine struct {
	sMutex  *sync.Mutex
	saved   []byte
	restore []byte
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
	if file.Path == "fail" || t.restore == nil {
		return nil, errors.New("fail")
	}

	return bytes.NewReader(t.restore), nil
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
func (t *TestEngine) SaveOptions() []engine.Option {
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

func (t *TestEngine) RestoreOptions() []engine.Option {
	return nil
}
func (t *TestEngine) ConfigureRestore(map[string]interface{}) error {
	return nil
}
