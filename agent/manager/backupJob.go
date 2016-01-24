package manager

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sethjback/gobble/agent/notification"
	"github.com/sethjback/gobble/agent/workers"
	"github.com/sethjback/gobble/engines"
	"github.com/sethjback/gobble/modifications"
	"github.com/sethjback/gobble/spec"
	"github.com/sethjback/gobble/util/log"
)

// BackupJob defines a backup job to run
// implements the Job interface
type BackupJob struct {
	ID                int
	Coordinator       *spec.Coordinator
	Start             string
	State             string
	Paramiters        spec.BackupParamiter
	Modifications     []modifications.Modification
	Engines           []engines.Backup
	FilesProccessed   int64
	cancel            chan struct{}
	MaxWorkers        int
	NotificationQueue *notification.Queue
}

// Configure verifies that all the necessary options have been provided for the modifications and eingines
// and confirms that the backup can in theory proceed
func (b *BackupJob) Configure(jobParams spec.BackupParamiter) error {
	b.Paramiters = jobParams

	mods, err := modifications.GetModifications(jobParams.Modifications, false)
	if err != nil {
		return err
	}

	b.Modifications = mods

	engines, err := engines.GetBackupEngines(jobParams.Engines)
	if err != nil {
		return err
	}

	b.Engines = engines

	b.State = "configured"

	return nil
}

// Status returns the current status of this backup
func (b *BackupJob) Status() map[string]interface{} {
	return map[string]interface{}{"start": b.Start, "processed": b.FilesProccessed, "status": b.State}
}

// Cancel will stop the bacup job. Any running workers will finish
func (b *BackupJob) Cancel() {
	log.Infof("backupJob", "Canceling backup job: %v", b.ID)
	b.State = "canceling"
	close(b.cancel)
}

// Run starts the backup job.
// Backups proceed as follows:
// It starts a goroutine to walk the JobPaths and return individual files over a channel
// It then starts a number of woker goroutines (MaxWorkers) to handle the actual work
// Each worker is handed the path channel and feeds off of incoming messages until it is closed
// Each worker is also handed a success channel for keeping track of progress
// It waits until all workes terminate, then closes the progress channel
// It then checks for any queued errors then sends a finished message over the finished channel
func (b *BackupJob) Run(finished chan<- bool) {

	log.Infof("backupJob", "running backupJob: %v", b.ID)
	log.Debugf("backupJob", "Backup Job: %v", *b)

	b.State = "running"
	b.cancel = make(chan struct{})
	b.Start = time.Now().String()
	b.NotificationQueue = notification.NewQueue(b.Coordinator, b.ID, *keyManager)
	paths, errc := buildBackupFileList(b.cancel, b.Paramiters.Paths)

	processed := make(chan int64)
	var wg sync.WaitGroup
	wg.Add(b.MaxWorkers)
	for i := 0; i < b.MaxWorkers; i++ {
		log.Debugf("backupJob", "Starting Worker: %v", i)
		go func() {
			backupWorker(b.cancel, paths, processed, b.Modifications, b.Engines, b.NotificationQueue.In, b.ID)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		log.Debugf("backupJob", "All Workers Done")
		close(processed)
	}()

	b.NotificationQueue.Run()

	for p := range processed {
		atomic.AddInt64(&b.FilesProccessed, p)
	}

	if err := <-errc; err != nil {
		//send errors to notifier
		//b.NotificationQueue.In <- map[string]interface{}{"Error": err}
	}

	b.State = "finishing notifications"

	f, _ := json.Marshal(&spec.JobUpdateRequest{spec.Complete})

	b.NotificationQueue.Finish(&notification.Notification{
		Endpoint: "/jobs/" + strconv.Itoa(b.ID) + "/complete",
		Payload:  f})

	//Wait for the notification queue to be done
	<-b.NotificationQueue.Finished

	b.State = "finished"

	finished <- true
}

// Builds the file list for the backup
// Listens for the cancel chanel to close to cancel walk
// Walks the file tree in JobPaths and sends any found file that isn't excluded on the return chan
// If there is an error, it sends the error on the error channel and returns
func buildBackupFileList(cancel <-chan struct{}, jobPaths []spec.BackupPath) (<-chan string, <-chan error) {

	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		log.Debug("backupJob", "file list routine started")
		defer close(paths)

		for _, jobPath := range jobPaths {

			log.Debugf("backupJob", "Walking filepath: %v", jobPath)

			errc <- filepath.Walk(jobPath.Path, func(path string, info os.FileInfo, err error) error {

				log.Debugf("backupJob", "Walk Found: %v", path)

				if err != nil {
					return err
				}

				if !info.Mode().IsRegular() {
					return nil
				}

				if shouldExclude(path, jobPath.Excludes) {
					return nil
				}

				select {
				case paths <- path:
				case <-cancel:
					log.Info("backupJob", "Walk Canceled")
					return errors.New("Walk Canceled")
				}
				return nil
			})
		}
	}()

	return paths, errc
}

// Function to pass incoming files on the paths chan to a waiting Backup worker
// Listens for close on cancel chan to stop
func backupWorker(cancel <-chan struct{}, paths <-chan string, progress chan<- int64, mods []modifications.Modification, engines []engines.Backup, nChan chan<- *notification.Notification, jobID int) {
	worker := &workers.Backup{mods, engines}
	for path := range paths {
		jf := worker.Do(path)
		payload, _ := json.Marshal(&spec.JobFileRequest{File: *jf})
		nChan <- &notification.Notification{Endpoint: "/jobs/" + strconv.Itoa(jobID) + "/files", Payload: payload}
		select {
		case progress <- 1:
		case <-cancel:
			return
		}
	}
}

// ShouldExclude takes a file path string and compares it with the requested exclusions
func shouldExclude(path string, excludes []string) bool {
	return false
}
