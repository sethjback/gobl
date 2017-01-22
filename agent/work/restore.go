package work

import (
	"io"

	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/modification"
	"github.com/sethjback/gobl/util/log"
)

// Restore defines the paramiters of the work to be done
type Restore struct {
	File          files.File
	Modifications []modification.Definition
	From          engine.Definition
	To            []engine.Definition
}

// Worker interface
func (r Restore) Do() interface{} {
	log.Debugf("restoreWorker", "Working on: %v", r.File.Signature.Path)
	jf := model.JobFile{}
	jf.File = r.File

	svrs, err := engine.BuildSavers([]engine.Definition{r.From})
	if err != nil {
		jf.State = StateErrors
		jf.Error = goblerr.New("unable to get from reader", ErrorRestoreEngines, "restore", err)
		return jf
	}

	//get the backup engine's Reader
	reader, err := svrs[0].Retrieve(r.File)
	if err != nil {
		jf.State = StateErrors
		jf.Error = goblerr.New("unable to get from reader", ErrorRestoreEngines, "restore", err)
		return jf
	}

	rers, err := engine.BuildRestorers(r.To)
	if err != nil {
		log.Infof("restore", "build to failed: %s", err)
		jf.Error = goblerr.New("unable to build to engines", ErrorRestoreEngines, "restore", err)
		jf.State = StateErrors

		return jf
	}

	eng, err := engine.NewRestoreEngine(r.File, rers...)
	if err != nil {
		jf.Error = goblerr.New("unable to build to engines", ErrorRestoreEngines, "restore", err)
		jf.State = StateErrors
		return jf
	}

	mods, err := modification.Build(r.Modifications, modification.Backward)
	if err != nil {
		jf.Error = goblerr.New("unable to build to modification pipeline", ErrorModifications, "restore", err)
		jf.State = StateErrors
		return jf
	}

	pipe := modification.Pipeline(reader, mods...)

	done := make(chan struct{})

	go func() {
		_, err := io.Copy(eng, pipe.Tail)
		if err != nil {
			pipe.Erroc <- err
		} else {
			eng.Finish()
			done <- struct{}{}
		}
	}()

	//Wait for an error or jobdone
	select {
	case err := <-pipe.Erroc:
		jf.Error = goblerr.New("file restore failed", ErrorRestore, "restore", err)
		jf.State = StateErrors
	case <-done:
		log.Debugf("restoreWorker", "Restore Done: %v", r.File.Path)
		jf.State = StateComplete
	}

	return jf
}

/*

// NewRestore returns a configured restore worker
// It requires that the modification slice be in the correct order for restoring
// and that the to and from engines be already configured
func NewRestore(mods []modification.Modifyer, to engine.Restorer, from engine.Saver) *Restore {
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
	pipe := modification.Pipeline(reader, errc, false, r.Modifications...)

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
*/
