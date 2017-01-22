// Package manager coordinates the runnign of jobs on the agent
// it mediates between the api and the actual job routines
package manager

import (
	"crypto/rsa"
	"errors"
	"runtime"
	"sync"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/spec"
	"github.com/sethjback/gobl/util"
	"github.com/sethjback/gobl/util/log"
)

// Error wraps additional data specific to the workers
type Error struct {
	Source  string `json:"source,omitempty"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *Error) Error() string { return e.Message }

// A Job is created and managed by the manager
// Can be a BackupJob or RestoreJob
type Job interface {
	Run(chan<- bool)
	Cancel()
	Status() map[string]interface{}
}

var active = make(map[int]Job)
var m sync.Mutex
var keyManager *keys.Manager
var conf *config.Config
var version util.Version

// Init configures the manager
func Init(c *config.Config) error {

	keyManager = &keys.Manager{PublicKeys: make(map[string]*rsa.PublicKey)}
	key, err := keys.OpenPrivateKey(c.Host.PrivateKeyPath)
	if err != nil {
		return err
	}

	keyManager.PrivateKey = key

	cKey, err := keys.OpenPublicKey(c.Coordinator.PublicKeyPath)
	if err != nil {
		return err
	}

	keyManager.PublicKeys[c.Coordinator.Address] = cKey

	conf = c

	return nil
}

// Status returns a list of information about the agent
func Status() map[string]interface{} {
	status := make(map[string]interface{})

	goR := runtime.NumGoroutine()
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	agentS := map[string]interface{}{"goRoutines": goR, "memory": memStat.Alloc, "jobs": len(active)}
	status["agent"] = agentS
	return status
}

// JobStatus returns the JobStatus of job ID provided
func JobStatus(id int) (interface{}, bool) {
	job, ok := getJob(id)
	if !ok {
		return &Error{"", "I can't find that job", 404}, false
	}

	return job.Status(), true
}

// NewRestore creates and starts a new restore job
func NewRestore(restoreRequest spec.RestoreJobRequest) error {

	restoreJob := new(RestoreJob)

	err := restoreJob.Configure(restoreRequest.Paramiters)

	if err != nil {
		return &Error{err.Error(), "Could not create Job", 400}
	}

	restoreJob.ID = restoreRequest.ID
	restoreJob.MaxWorkers = 3
	restoreJob.Coordinator = &conf.Coordinator
	restoreJob.Paramiters = restoreRequest.Paramiters

	finishedChan := make(chan bool)
	addJob(restoreJob.ID, restoreJob)
	go restoreJob.Run(finishedChan)
	go wait(restoreJob.ID, finishedChan)

	return nil
}

// NewBackup creates a new Job worker and starts
func NewBackup(backupRequest spec.BackupJobRequest) error {

	backupJob := new(BackupJob)

	if len(backupRequest.Paramiters.Paths) == 0 || len(backupRequest.Paramiters.Modifications) == 0 || len(backupRequest.Paramiters.Engines) == 0 {
		return errors.New("Invalid paramiters provided")
	}

	err := backupJob.Configure(backupRequest.Paramiters)
	if err != nil {
		return &Error{err.Error(), "Could Not Create Job", 400}
	}

	backupJob.Coordinator = &conf.Coordinator
	backupJob.ID = backupRequest.ID
	// TODO: make configurable
	backupJob.MaxWorkers = 3

	finishedChan := make(chan bool)
	addJob(backupJob.ID, backupJob)
	go backupJob.Run(finishedChan)
	go wait(backupJob.ID, finishedChan)

	return nil
}

// Cancel stopas a currently running Job
func Cancel(id int) (bool, error) {

	m.Lock()
	defer m.Unlock()
	if job, ok := active[id]; ok {
		job.Cancel()
		return true, nil
	}

	return false, &Error{"", "I can't find that job", 404}
}

func wait(on int, over chan bool) {
	<-over
	log.Infof("manager", "Job Finished: %v", on)
	removeJob(on)
}

func getJobs() []Job {
	var rCopy []Job
	for _, job := range active {
		rCopy = append(rCopy, job)
	}
	m.Unlock()
	return rCopy
}

func removeJob(id int) {
	m.Lock()
	delete(active, id)
	m.Unlock()
}

func addJob(id int, job Job) {
	m.Lock()
	active[id] = job
	m.Unlock()
}

func getJob(id int) (Job, bool) {
	m.Lock()
	val, ok := active[id]
	m.Unlock()
	return val, ok
}

// GetKey returns the agent's public key as base64.urlencoded string
func GetKey() (string, error) {
	return keyManager.PublicKey()
}

// VerifySignature checks the incoming request agains
func VerifySignature(signed []byte, signature string) error {
	log.Debugf("manager", "verify sig: %v", conf.Coordinator)
	key, err := keyManager.KeyForHost(conf.Coordinator.Address)
	if err != nil {
		return errors.New("Cannot find coordinator key")
	}
	return keys.VerifySignature(key, signed, signature)
}
