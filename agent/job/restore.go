package job

import (
	"sync"
	"time"

	"github.com/sethjback/gobl/agent/notification"
	"github.com/sethjback/gobl/agent/work"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/util/log"
	"github.com/sethjback/gowork"
)

type Restore struct {
	stateM      *sync.Mutex
	Job         model.Job
	Coordinator config.Coordinator
	cancel      chan struct{}
	MaxWorkers  int
	Notifier    notification.Notifier
}

// Status for the jobber interface
func (r *Restore) Status() model.JobMeta {
	r.stateM.Lock()
	jm := *r.Job.Meta
	r.stateM.Unlock()
	return jm
}

// Cancel for jobber interface
func (r *Restore) Cancel() {
	log.Infof("job", "Cancel restore: %v", r.Job.ID)
	r.stateM.Lock()
	r.Job.Meta.State = StateCanceling
	close(r.cancel)
	r.stateM.Unlock()
}

func (r *Restore) SetState(state string) {
	r.stateM.Lock()
	r.Job.Meta.State = state
	r.stateM.Unlock()
}

func (r *Restore) GetState() string {
	r.stateM.Lock()
	s := r.Job.Meta.State
	r.stateM.Unlock()
	return s
}

func (r *Restore) addTotal(num int) {
	r.stateM.Lock()
	r.Job.Meta.Total += num
	r.stateM.Unlock()
}

func (r *Restore) addComplete(num int) {
	r.stateM.Lock()
	r.Job.Meta.Complete += num
	r.stateM.Unlock()
}

func (r *Restore) Run(finished chan<- string) {
	log.Infof("restoreJob", "running restorJob: %v", r.Job.ID)
	log.Debugf("restoreJob", "Restore Job: %v", *r)

	r.SetState(StateRunning)
	r.cancel = make(chan struct{})
	r.Job.Meta.Start = time.Now()

	q := gowork.NewQueue(100, r.MaxWorkers)
	q.Start(r.MaxWorkers)

	go func() {
		r.addTotal(len(r.Job.Definition.Files))
		for _, f := range r.Job.Definition.Files {
			q.AddWork(work.Restore{File: f, From: r.Job.Definition.From, To: r.Job.Definition.To, Modifications: r.Job.Definition.Modifications})
		}
		q.Finish()
	}()

	done := make(chan struct{})

	go func() {
		processedFiles := 0
		for result := range q.Results() {
			//send to notification Q
			jf := result.(model.JobFile)
			r.Notifier.Send(&JobNotification{JF: &jf, dest: r.Coordinator.Address})
			processedFiles++
			if processedFiles > 10 {
				r.addComplete(processedFiles)
				processedFiles = 0
			}
			r.addComplete(processedFiles)
		}
		log.Debug("restoreJob", "q results closed")
		close(done)
	}()

	select {
	case <-r.cancel:
		q.Abort()
		// wait for graceful shutdown
		<-done

	case <-done:
		//finished!
	}

	// notify our manager that we are done
	finished <- r.Job.ID
}
