package job

import (
	"sync"

	"google.golang.org/grpc"

	"github.com/sethjback/gobl/agent/grpcclient"
	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/modification"
)

const (
	ErrorCreateJob = "CreateJobFailed"
	ErrorFindJob   = "FindJobFailed"
)

var jobs = make(map[string]Jobber)
var jobMutex sync.Mutex
var finish = make(chan string)
var grpcConn *grpc.ClientConn

// Jobber
type Jobber interface {
	Run(done chan<- string)
	Cancel()
	Status() model.JobMeta
}

func Init() {
	go waiter()
}

func Shutdown() {

}

func waiter() {
	for id := range finish {
		removeJob(id)
		if pendingJobs() == 0 {
			grpcclient.CloseClient()
		}
	}
}

func removeJob(id string) {
	jobMutex.Lock()
	delete(jobs, id)
	jobMutex.Unlock()
}

func addJob(id string, job Jobber) {
	jobMutex.Lock()
	jobs[id] = job
	jobMutex.Unlock()
}

func pendingJobs() int {
	jobMutex.Lock()
	defer jobMutex.Unlock()
	return len(jobs)
}

func State(id string) (model.JobMeta, error) {
	var status model.JobMeta
	var err error
	jobMutex.Lock()
	j, ok := jobs[id]
	if ok {
		status = j.Status()
	} else {
		err = goblerr.New("I was unable to find that Job", ErrorFindJob, nil)
	}
	jobMutex.Unlock()
	return status, err
}

func Cancel(id string) (err error) {
	jobMutex.Lock()
	if j, ok := jobs[id]; ok {
		j.Cancel()
	} else {
		err = goblerr.New("I was unable to find that Job", ErrorFindJob, nil)
	}
	jobMutex.Unlock()
	return
}

func NewRestore(job model.Job) error {
	r := &restore{
		stateM:     &sync.Mutex{},
		job:        job,
		MaxWorkers: 3,
	}
	r.job.Meta = &model.JobMeta{}

	_, err := engine.BuildSavers([]engine.Definition{*r.job.Definition.From})
	if err != nil {
		return err
	}

	_, err = engine.BuildRestorers(r.job.Definition.To)
	if err != nil {
		return err
	}

	_, err = modification.Build(r.job.Definition.Modifications, modification.Backward)
	if err != nil {
		return err
	}

	addJob(job.ID, r)

	r.coordClient, err = grpcclient.Client()
	if err != nil {
		removeJob(job.ID)
		return err
	}

	go r.Run(finish)

	return nil
}

func NewBackup(job model.Job) error {
	b := &backup{
		stateM:     &sync.Mutex{},
		job:        job,
		MaxWorkers: 3,
	}
	b.job.Meta = &model.JobMeta{}

	_, err := engine.BuildSavers(b.job.Definition.To)
	if err != nil {
		return err
	}

	_, err = modification.Build(b.job.Definition.Modifications, modification.Forward)
	if err != nil {
		return err
	}

	addJob(job.ID, b)

	b.coordClient, err = grpcclient.Client()
	if err != nil {
		removeJob(job.ID)
		return err
	}

	go b.Run(finish)
	return nil
}
