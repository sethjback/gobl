package job

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sethjback/gobl/agent/coordinator"
	"github.com/sethjback/gobl/agent/notification"
	"github.com/sethjback/gobl/agent/work"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gowork"
)

type Backup struct {
	stateM      *sync.Mutex
	Job         model.Job
	Coordinator *coordinator.Coordinator
	cancel      chan struct{}
	MaxWorkers  int
	Notifier    notification.Notifier
}

func NewBackup(job model.Job, coordinator *coordinator.Coordinator, notifier notification.Notifier) (*Backup, error) {
	job.Meta = &model.JobMeta{}
	return &Backup{
		stateM:      &sync.Mutex{},
		Job:         job,
		Coordinator: coordinator,
		Notifier:    notifier,
		MaxWorkers:  3,
	}, nil
}

// Status for the jobber interface
func (b *Backup) Status() model.JobMeta {
	b.stateM.Lock()
	jm := *b.Job.Meta
	b.stateM.Unlock()
	return jm
}

// Cancel for jobber interface
func (b *Backup) Cancel() {
	b.stateM.Lock()
	b.Job.Meta.State = model.StateCanceling
	close(b.cancel)
	b.stateM.Unlock()
}

func (b *Backup) SetState(state string) {
	b.stateM.Lock()
	b.Job.Meta.State = state
	b.stateM.Unlock()
}

func (b *Backup) GetState() string {
	b.stateM.Lock()
	s := b.Job.Meta.State
	b.stateM.Unlock()
	return s
}

func (b *Backup) addTotal(num int) {
	b.stateM.Lock()
	b.Job.Meta.Total += num
	b.stateM.Unlock()
}

func (b *Backup) addComplete(num int) {
	b.stateM.Lock()
	b.Job.Meta.Complete += num
	b.stateM.Unlock()
}

// Run for jobber interface
func (b *Backup) Run(finished chan<- string) {
	b.SetState(model.StateRunning)
	b.cancel = make(chan struct{})
	b.Job.Meta.Start = time.Now()

	paths, errc := buildBackupFileList(b.cancel, b.Job.Definition.Paths)

	q := gowork.NewQueue(100, b.MaxWorkers)
	q.Start(b.MaxWorkers)

	go func() {
		totalFiles := 0
		for path := range paths {
			q.AddWork(work.Backup{File: path, Modifications: b.Job.Definition.Modifications, Engines: b.Job.Definition.To})
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

	done := make(chan struct{})

	go func() {
		processedFiles := 0
		for result := range q.Results() {
			//send to notification Q
			jf := result.(model.JobFile)
			b.Notifier.Send(&JobNotification{JF: &jf, host: b.Coordinator.Address, path: "/jobs/" + b.Job.ID + "/files"})
			processedFiles++
			if processedFiles > 10 {
				b.addComplete(processedFiles)
				processedFiles = 0
			}
			b.addComplete(processedFiles)
		}
		close(done)
	}()

	// file walk error
	if err := <-errc; err != nil {
		//send errors to notifier
		//b.NotificationQueue.In <- map[string]interface{}{"Error": err}
	}

	select {
	case <-b.cancel:
		q.Abort()
		// wait for graceful shutdown
		<-done

	case <-done:
		//finished!
	}

	b.Notifier.Send(&JobNotification{host: b.Coordinator.Address, path: "/jobs/" + b.Job.ID + "/complete"})

	// notify our manager that we are done
	finished <- b.Job.ID
}

// Builds the file list for the backup
// Listens for the cancel chanel to close to cancel walk
// Walks the file tree in JobPaths and sends any found file that isn't excluded on the return chan
// If there is an error, it sends the error on the error channel and returns
func buildBackupFileList(cancel <-chan struct{}, paths []model.Path) (<-chan string, <-chan error) {

	files := make(chan string)
	errc := make(chan error, 1)

	go func() {

		for _, path := range paths {

			errc <- filepath.Walk(path.Root, func(filePath string, info os.FileInfo, err error) error {

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
		}

		close(files)
		close(errc)
	}()

	return files, errc
}

// ShouldExclude takes a file path string and compares it with the requested exclusions
func shouldExclude(path string, excludes []string) bool {
	// TODO: implement
	return false
}
