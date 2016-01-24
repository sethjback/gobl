package manager

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sethjback/gobble/agent/notification"
	"github.com/sethjback/gobble/agent/workers"
	"github.com/sethjback/gobble/engines"
	"github.com/sethjback/gobble/files"
	"github.com/sethjback/gobble/modifications"
	"github.com/sethjback/gobble/spec"
	"github.com/sethjback/gobble/util/log"
)

// RestoreJob defines a restore job to be run
// implements the Job interface
type RestoreJob struct {
	ID                int
	coordinator       *spec.Coordinator
	Start             string
	State             string
	Paramiters        spec.RestoreParamiter
	Modifications     []modifications.Modification
	From              *engines.Backup
	To                *engines.Restore
	cancel            chan struct{}
	MaxWorkers        int
	FilesProccessed   int64
	NotificationQueue *notification.Queue
}

// Configure the restore job
func (r *RestoreJob) Configure(rp spec.RestoreParamiter) error {

	//configure modificaitons and reverse
	mods, err := modifications.GetModifications(rp.BackupParamiters.Modifications, true)
	if err != nil {
		return err
	}

	r.Modifications = mods

	//Confirm the Engine we want to restore from was part of the job and configure
	for _, e := range rp.BackupParamiters.Engines {
		if e.Name == rp.From.Name {
			be, err := engines.GetBackupEngines([]engines.Definition{e})
			if err != nil {
				return errors.New("Could not create restore from engine instance: " + err.Error())
			}
			r.From = &be[0]
		}
	}

	if r.From == nil {
		return errors.New("Engine to restore from must have been part of the original backup")
	}

	re, err := engines.GetRestoreEngine(rp.To)
	if err != nil {
		return err
	}
	r.To = &re

	return nil
}

// Definition returns the job definitions
func (r *RestoreJob) Definition() map[string]interface{} {
	return map[string]interface{}{
		"restoreDefinition": r.Paramiters}
}

// Status returns the job status
func (r *RestoreJob) Status() map[string]interface{} {
	return map[string]interface{}{"start": r.Start, "processed": r.FilesProccessed}
}

// Cancel the restore job
func (r *RestoreJob) Cancel() {
	close(r.cancel)
}

// Run the defined restore operation
func (r *RestoreJob) Run(finished chan<- bool) {

	log.Infof("restoreJob", "running restoreJob: %v", r.ID)
	log.Debugf("restoreJob", "Restore Job: %v", *r)

	r.State = "running"
	r.cancel = make(chan struct{})
	r.Start = time.Now().String()
	r.NotificationQueue = notification.NewQueue(r.coordinator, r.ID, *keyManager)

	sigsc := make(chan files.Signature)
	processed := make(chan int64)

	go func() {
		for _, sig := range r.Paramiters.FileSignatures {
			sigsc <- sig
		}
		close(sigsc)
	}()

	var wg sync.WaitGroup
	wg.Add(r.MaxWorkers)
	for i := 0; i < r.MaxWorkers; i++ {
		log.Debugf("restoreJob", "Starting Worker: %v", i)
		go func() {
			restoreWorker(r.cancel, sigsc, processed, r.Modifications, *r.To, *r.From, r.NotificationQueue.In, r.ID)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		log.Debugf("restoreJob", "All Workers Done")
		close(processed)
	}()

	r.NotificationQueue.Run()

	for p := range processed {
		atomic.AddInt64(&r.FilesProccessed, p)
	}

	r.State = "finishing notifications"

	f, _ := json.Marshal(&spec.JobUpdateRequest{spec.Complete})

	r.NotificationQueue.Finish(&notification.Notification{
		Endpoint: "/jobs/" + strconv.Itoa(r.ID) + "/complete",
		Payload:  f})

	//Wait for the notification queue to be done
	<-r.NotificationQueue.Finished

	r.State = "finished"

	finished <- true
}

func restoreWorker(cancel <-chan struct{}, sigsc <-chan files.Signature, progress chan<- int64, mods []modifications.Modification, to engines.Restore, from engines.Backup, nChan chan<- *notification.Notification, jobID int) {
	worker := &workers.Restore{mods, from, to}
	for sig := range sigsc {
		jf := worker.Do(&sig)
		payload, _ := json.Marshal(&spec.JobFileRequest{File: *jf})
		nChan <- &notification.Notification{Endpoint: "/jobs/" + strconv.Itoa(jobID) + "/files", Payload: payload}
		select {
		case progress <- 1:
		case <-cancel:
			return
		}
	}
}
