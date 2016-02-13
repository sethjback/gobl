package workers

import (
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/kalafut/imohash"
	"github.com/sethjback/gobl/engines"
	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/modifications"
	"github.com/sethjback/gobl/spec"
)

// Backup defines the paramiters of the work to be done
type Backup struct {
	Modifications []modifications.Modification
	Engines       []engines.Backup
}

// NewBackup Creates a new Backup worker
func NewBackup(mods []modifications.Modification, engines []engines.Backup) *Backup {
	return &Backup{mods, engines}
}

// Do back up the given File using Modifications to Engines
func (b *Backup) Do(file string) *spec.JobFile {
	jf := &spec.JobFile{}

	//build the signature
	dir, name := filepath.Split(file)

	var ms []string
	for _, mod := range b.Modifications {
		ms = append(ms, mod.Name())
	}

	fileSig := files.Signature{
		Name:          name,
		Path:          dir,
		Modifications: ms}

	imoHash := imohash.NewCustom(24*1024, 3*1024*1024)
	fileHash, err := imoHash.SumFile(file)
	if err != nil {
		//log error
		//fmt.Println(err)
		jf.Message = err.Error()
		jf.State = spec.Errors
		jf.Signature = fileSig

		return jf
	}

	fileSig.Hash = hex.EncodeToString(fileHash[:])

	jf.Signature = fileSig

	var saveEngines []*engines.Backup

	for _, e := range b.Engines {
		store, err := e.ShouldBackup(fileSig)
		if err != nil {
			//Deal with the error!
			//fmt.Println(err)
			jf.Message = err.Error()
			jf.State = spec.Errors
			return jf
		}

		if store {
			saveEngines = append(saveEngines, &e)
		}
	}

	if len(saveEngines) == 0 {
		//Nothing to do!
		//fmt.Println("no save engines, Job Done")
		jf.State = spec.Complete
		jf.Message = "no save engines, Job Done"
		return jf
	}

	fileHandle, err := os.Open(file)
	if err != nil {
		//handle Error
		jf.Message = err.Error()
		jf.State = spec.Errors
		return jf
	}
	defer fileHandle.Close()

	errc := make(chan error)
	done := make(chan bool)

	//configure the Engines
	ew := &engines.Writer{}
	ew.Engines = make([]*io.PipeWriter, len(b.Engines))
	defer ew.Close()

	for i := 0; i < len(b.Engines); i++ {
		r, w := io.Pipe()
		go b.Engines[i].Backup(r, fileSig, errc)
		ew.Engines[i] = w
	}

	//setup the pipeline
	pipe := modifications.NewPipeline(fileHandle, errc, true, b.Modifications...)

	//copy the data
	go func() {
		if _, err := io.Copy(ew, pipe.Tail); err != nil {
			errc <- err
			return
		}
		done <- true
	}()

	//Wait for an error or jobdone
	select {
	case err := <-errc:
		jf.Message = err.Error()
		jf.State = spec.Errors
		return jf
	case <-done:
		jf.State = spec.Complete
		return jf
	}
}
