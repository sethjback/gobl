package work

import (
	"encoding/hex"
	"io"
	"os"

	"github.com/kalafut/imohash"
	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/modification"
	"github.com/sethjback/gobl/util/log"
)

// Backup defines the paramiters of the work to be done
type Backup struct {
	File          string
	Modifications []modification.Definition
	Engines       []engine.Definition
}

// The Work interface from the worker package
func (b Backup) Do() interface{} {
	jf := model.JobFile{}

	jf.File.Signature = files.NewSignature(b.File, b.Modifications)

	imoHash := imohash.NewCustom(24*1024, 3*1024*1024)
	fileHash, err := imoHash.SumFile(b.File)
	if err != nil {
		log.Infof("backupWork", "(%s) hash failed: %s", b.File, err)
		jf.Error = goblerr.New("unable to hash file", ErrorFileHash, "backup", err)
		jf.State = StateErrors
		return jf
	}

	jf.File.Signature.Hash = hex.EncodeToString(fileHash[:])

	svrs, err := engine.BuildSavers(b.Engines)
	if err != nil {
		log.Infof("backupWork", "build savers failed: %s", err)
		jf.Error = goblerr.New("unable bulid save engines", ErrorSaveEngines, "backup", err)
		jf.State = StateErrors

		return jf
	}

	eng, saveNeeded, err := engine.NewBackupEngine(jf.File, svrs...)
	if err != nil {
		jf.Error = goblerr.New("unable bulid save engines", ErrorSaveEngines, "backup", err)
		jf.State = StateErrors
		return jf
	}

	// no savers
	if !saveNeeded {
		jf.State = StateSkipped
		return jf
	}

	fileHandle, err := os.Open(b.File)
	if err != nil {
		jf.Error = goblerr.New("unable to open file", ErrorFileOps, "backup", err)
		jf.State = StateErrors
		return jf
	}
	defer fileHandle.Close()

	mods, err := modification.Build(b.Modifications, modification.Forward)
	if err != nil {
		jf.Error = goblerr.New("unable bulid modifications", ErrorModifications, "backup", err)
		jf.State = StateErrors
		return jf
	}

	pipe := modification.Pipeline(fileHandle, mods...)

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
		jf.Error = goblerr.New("file save failed", ErrorSave, "backup", err)
		jf.State = StateErrors
	case <-done:
		jf.State = StateComplete
	}

	return jf
}
