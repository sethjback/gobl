package job

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sethjback/gobl/agent/work"
	pb "github.com/sethjback/gobl/goblgrpc"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gowork"
)

type backup struct {
	stateM      *sync.Mutex
	job         model.Job
	cancel      chan struct{}
	MaxWorkers  int
	coordClient pb.CoordinatorClient
}

// Status for the jobber interface
func (b *backup) Status() model.JobMeta {
	b.stateM.Lock()
	jm := *b.job.Meta
	b.stateM.Unlock()
	return jm
}

// Cancel for jobber interface
func (b *backup) Cancel() {
	b.stateM.Lock()
	b.job.Meta.State = model.StateCanceling
	close(b.cancel)
	b.stateM.Unlock()
}

func (b *backup) SetState(state string) {
	b.stateM.Lock()
	b.job.Meta.State = state
	b.stateM.Unlock()
}

func (b *backup) GetState() string {
	b.stateM.Lock()
	s := b.job.Meta.State
	b.stateM.Unlock()
	return s
}

func (b *backup) addTotal(num int) {
	b.stateM.Lock()
	b.job.Meta.Total += num
	b.stateM.Unlock()
}

func (b *backup) addComplete(num int) {
	b.stateM.Lock()
	b.job.Meta.Complete += num
	b.stateM.Unlock()
}

// Run for jobber interface
func (b *backup) Run(finished chan<- string) {
	b.SetState(model.StateRunning)
	b.cancel = make(chan struct{})
	b.job.Meta.Start = time.Now()

	completed := make(chan *model.JobFile, 20)
	errc := make(chan error)
	wg := &sync.WaitGroup{}

	paths := buildBackupFileList(b.cancel, b.job.Definition.Paths, errc)

	q := gowork.NewQueue(100, b.MaxWorkers)
	q.Start(b.MaxWorkers)

	wg.Add(1)
	// range over the walked files and send them to the work queue
	go func() {
		defer wg.Done()
		totalFiles := 0
		for path := range paths {
			q.AddWork(work.Backup{File: path, Modifications: b.job.Definition.Modifications, Engines: b.job.Definition.To})
			totalFiles++
			if totalFiles > 10 {
				b.addTotal(totalFiles)
				totalFiles = 0
			}
		}
		b.addTotal(totalFiles)
		// close the input channel into the work queue
		q.Finish()
	}()

	wg.Add(1)
	// range over the completed work and send to notification queue
	go func() {
		defer wg.Done()
		defer close(completed)
		processedFiles := 0
		for result := range q.Results() {
			//send to notification Q
			jf := result.(model.JobFile)
			completed <- &jf
			processedFiles++
			if processedFiles > 10 {
				b.addComplete(processedFiles)
				processedFiles = 0
			}
		}
		b.addComplete(processedFiles)
	}()

	wg.Add(1)
	// range over completed files and send to coordinator
	go func() {
		defer wg.Done()
		fileClient, err := b.coordClient.File(context.Background())
		if err != nil {
			fmt.Printf("Error file client: %+v\n", err)
			errc <- err
			return
		}
		for file := range completed {
			fr := &pb.FileRequest{
				JobId: b.job.ID,
				File: &pb.File{
					Signature: &pb.Signature{
						Path:          file.File.Signature.Path,
						Hash:          file.File.Signature.Hash,
						Modifications: file.File.Signature.Modifications,
					},
					Meta: &pb.Meta{
						Mode: file.File.Meta.Mode,
						Uid:  int32(file.File.Meta.UID),
						Gid:  int32(file.File.Meta.GID),
					},
				},
			}

			if file.State == work.StateErrors {
				fr.State = pb.State_FAILED
				fr.Message = file.Error
			} else {
				fr.State = pb.State_FINISHED
			}

			fileClient.Send(fr)
		}

		_, err = fileClient.CloseAndRecv()
		if err != nil {
			fmt.Printf("Error close and receive: %+v\n", err)
			errc <- err
		}

	}()

	// wait for all routines to finish then close the error channel
	go func() {
		wg.Wait()
		close(errc)
	}()

	// once the errors are processed close done
	done := make(chan struct{})
	go func() {
		defer close(done)
		for err := range errc {
			fmt.Printf("%s : error: %s\n", b.job.ID, err)
		}
	}()

	select {
	case <-b.cancel:
		q.Abort()
		// wait for graceful finish
		<-done
	case <-done:
		//finished!
	}

	// notify our manager that we are done
	finished <- b.job.ID
}

// Builds the file list for the backup
// Listens for the cancel chanel to close to cancel walk
// Walks the file tree in JobPaths and sends any found file that isn't excluded on the return chan
// If there is an error, it sends the error on the error channel and returns
func buildBackupFileList(cancel <-chan struct{}, paths []model.Path, errc chan<- error) <-chan string {

	files := make(chan string)

	go func() {

		for _, path := range paths {

			err := filepath.Walk(path.Root, func(filePath string, info os.FileInfo, err error) error {

				if err != nil {
					return err
				}

				if !info.Mode().IsRegular() {
					return nil
				}

				if shouldExclude(filePath, path.Excludes) {
					return nil
				}

				select {
				case files <- filePath:
				case <-cancel:
					return errors.New("Walk Canceled")
				}
				return nil
			})

			if err != nil {
				errc <- err
			}
		}

		close(files)
	}()

	return files
}

// ShouldExclude takes a file path string and compares it with the requested exclusions
func shouldExclude(path string, excludes []string) bool {
	// TODO: implement
	return false
}
