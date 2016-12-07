package engine

import (
	"io"

	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/goblerr"
)

const (
	ErrorFileCheck = "FileCheckFailed"
)

// Engine handles routing incoming data through to all the configured Savers and Restorers
// It implements the io.Writer interface so it can be passed to things like io.Copy()
type Engine interface {
	// Implements io.Writer interface to save the data
	io.Writer
	// ErrorChan returns the channel that all writers will send errors over
	ErrorChan() <-chan error
	// Finish must be called to close the writers
	Finish()
}

// backupEngine is the type that implements the Engine interface for backups
type backupEngine struct {
	savers []Saver
	pipes  []*io.PipeWriter
	errc   chan error
}

// NewBackupEngine returns an engine configured
func NewBackupEngine(signature files.Signature, savers ...Saver) (Engine, error) {
	e := &backupEngine{savers: savers}
	e.errc = make(chan error)
	for i := 0; i < len(e.savers); i++ {
		ok, err := e.savers[i].ShouldSave(signature)
		if err != nil {
			return nil, goblerr.New("Saver error on file check", ErrorFileCheck, nil, e.savers[i].Name()+" errored on ShouldBackup for "+signature.Name)
		}
		if ok {
			r, w := io.Pipe()
			go e.savers[i].Save(r, signature, e.errc)
			e.pipes = append(e.pipes, w)
		}
	}

	return e, nil
}

func (b *backupEngine) ErrorChan() <-chan error {
	return b.errc
}

// Write is just io.MultiWriter
func (b *backupEngine) Write(p []byte) (n int, err error) {
	for _, w := range b.pipes {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

// Closes all the engine pipes
func (b *backupEngine) Finish() {
	for _, w := range b.pipes {
		w.Close()
	}
}

// restoreEngine is the type that implements the Engine interface to restore
type restoreEngine struct {
	to    []Restorer
	pipes []*io.PipeWriter
	errc  chan error
}

func NewRestoreEngine(signature files.Signature, to ...Restorer) (Engine, error) {
	e := &restoreEngine{to: to}
	e.errc = make(chan error)
	for i := 0; i < len(e.to); i++ {
		ok, err := e.to[i].ShouldRestore(signature)
		if err != nil {
			return nil, goblerr.New("Restorer error on file check", ErrorFileCheck, nil, e.to[i].Name()+" errored on ShouldRestore for "+signature.Name)
		}
		if ok {
			r, w := io.Pipe()
			go e.to[i].Restore(r, signature, e.errc)
			e.pipes = append(e.pipes, w)
		}
	}
	return e, nil
}

func (r *restoreEngine) ErrorChan() <-chan error {
	return r.errc
}

func (r *restoreEngine) Write(p []byte) (n int, err error) {
	for _, w := range r.pipes {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

func (r *restoreEngine) Finish() {
	for _, w := range r.pipes {
		w.Close()
	}
}
