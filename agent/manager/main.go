// Package manager coordinates the runnign of jobs on the agent
// it mediates between the api and the actual job routines
package manager

import (
	"runtime"
	"sync"

	"github.com/sethjback/gobl/agent/job"
	"github.com/sethjback/gobl/agent/notification"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/util/log"
)

const (
	ErrorCreateJob = "CreateJobFailed"
	ErrorFindJob   = "FindJobFailed"
	ErrorReadKey   = "OpenKeyFileFailed"
	ErrorDecodeKey = "DecodeKeyFailed"
)

var active = make(map[string]job.Jobber)
var jobMutex sync.Mutex
var stateMutex sync.Mutex
var conf *config.Config
var notifier notification.Notifier
var finish chan string
var running bool

// Init configures the manager
func Init(c *config.Config) error {
	conf = c

	aKey, err := keys.OpenPrivateKey(conf.Server.PrivateKey)
	if err != nil {
		return err
	}

	notifier = notification.New(&notification.Config{MaxWorkers: 3, MaxDepth: 6}, aKey)
	notifier.Start()

	finish = make(chan string)
	running = true

	go waiter()

	return nil
}

func Shutdown() {
	stateMutex.Lock()
	running = false
	stateMutex.Unlock()

	jobMutex.Lock()
	for _, j := range active {
		j.Cancel()
	}
	jobMutex.Unlock()

	// will close once all jobs return
	<-finish
}

// Status returns a list of information about the agent
func Status() map[string]interface{} {
	status := make(map[string]interface{})

	goR := runtime.NumGoroutine()
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	agentS := map[string]interface{}{"goRoutines": goR, "memory": memStat.Alloc, "jobs": jobIds()}
	status["agent"] = agentS
	return status
}

// NewRestore creates and starts a new restore job
func NewRestore(restoreJob model.Job) goblerr.Error {
	r, err := job.NewRestore(restoreJob, conf.Coordinator, notifier)
	if err != nil {
		return goblerr.New("Unable to create job", ErrorCreateJob, "manager", err.Error())
	}

	addJob(restoreJob.ID, r)
	go r.Run(finish)

	return nil
}

// NewBackup creates a new Job worker and starts
func NewBackup(backupJob model.Job) goblerr.Error {
	b, err := job.NewBackup(backupJob, conf.Coordinator, notifier)
	if err != nil {
		return goblerr.New("Unable to create job", ErrorCreateJob, "manager", err.Error())
	}

	addJob(backupJob.ID, b)
	go b.Run(finish)

	return nil
}

// Cancel stops a currently running Job
func Cancel(id string) goblerr.Error {
	jobMutex.Lock()
	defer jobMutex.Unlock()
	if job, ok := active[id]; ok {
		job.Cancel()
		return nil
	}

	return goblerr.New("I was unable to find that Job", ErrorFindJob, "manager", "")
}

func waiter() {
	defer close(finish)
	for id := range finish {
		log.Infof("manager", "Job Finished: %v", id)
		removeJob(id)
		if stopped() && len(jobIds()) == 0 {
			return
		}
	}
}

func stopped() bool {
	var state bool
	stateMutex.Lock()
	state = running
	stateMutex.Unlock()
	return !state
}

func removeJob(id string) {
	jobMutex.Lock()
	delete(active, id)
	jobMutex.Unlock()
}

func addJob(id string, job job.Jobber) {
	jobMutex.Lock()
	active[id] = job
	jobMutex.Unlock()
}

func JobStatus(id string) (model.JobMeta, goblerr.Error) {
	var status model.JobMeta
	var err goblerr.Error
	jobMutex.Lock()
	j, ok := active[id]
	if ok {
		status = j.Status()
	} else {
		err = goblerr.New("I was unable to find that Job", ErrorFindJob, "manager", "")
	}
	jobMutex.Unlock()
	return status, err
}

func jobIds() []string {
	ids := []string{}
	jobMutex.Lock()
	for id, _ := range active {
		ids = append(ids, id)
	}
	jobMutex.Unlock()
	return ids
}

func Key() (string, goblerr.Error) {
	aKey, err := keys.OpenPrivateKey(conf.Server.PrivateKey)
	if err != nil {
		return "", goblerr.New("Unable to open key", ErrorReadKey, "manager", err.Error())
	}

	pks, err := keys.PublicKey(aKey)
	if err != nil {
		return "", goblerr.New("Unable to open key", ErrorDecodeKey, "manager", err.Error())
	}

	return pks, nil
}
