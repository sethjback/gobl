package workers

import (
	"bufio"
	"io"

	"github.com/sethjback/gobl/engines"
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/modifications"
	"github.com/sethjback/gobl/spec"
	"github.com/sethjback/gobl/util/log"
)

// Restore defines the paramiters of the work to be done
type Restore struct {
	Modifications []modifications.Modification
	From          engines.Backup
	To            engines.Restore
}

// NewRestore returns a configured restore worker
// It requires that the modification slice be in the correct order for restoring
// and that the to and from engines be already configured
func NewRestore(mods []modifications.Modification, to engines.Restore, from engines.Backup) *Restore {
	return &Restore{mods, from, to}
}

// Do restore the given file signature from the backup engine to the restore engine
func (r *Restore) Do(fileSig *files.Signature) *spec.JobFile {
	log.Debugf("restoreWorker", "Working on: %v", fileSig.Name)
	jf := &spec.JobFile{}
	jf.Signature = *fileSig

	errc := make(chan error)
	done := make(chan bool)
	defer close(errc)
	defer close(done)

	//get the backup engine's Reader
	reader, err := r.From.Retrieve(*fileSig)
	if err != nil {
		jf.State = spec.Errors
		jf.Message = err.Error()
		return jf
	}

	//setup the restore pipe
	pipeR, pipeW := io.Pipe()

	//configure the restore point
	go r.To.Restore(pipeR, *fileSig, errc)

	//setup the decode pipeline
	pipe := modifications.NewPipeline(bufio.NewReader(reader), errc, false, r.Modifications...)

	//copy the data
	go func() {
		if _, err := io.Copy(pipeW, pipe.Tail); err != nil {
			errc <- err
			return
		}
		done <- true
	}()

	//Wait for an error or jobdone
	select {
	case err := <-errc:
		jf.State = spec.Errors
		jf.Message = err.Error()
	case <-done:
		log.Debugf("restoreWorker", "Restore Done: %v", fileSig.Name)
		jf.State = spec.Complete
	}

	return jf
}
