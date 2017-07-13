package job

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/sethjback/gobl/agent/work"
	"github.com/sethjback/gobl/files"
	pb "github.com/sethjback/gobl/goblgrpc"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gowork"
)

type restore struct {
	stateM      *sync.Mutex
	job         model.Job
	cancel      chan struct{}
	MaxWorkers  int
	coordClient pb.CoordinatorClient
}

// Status for the jobber interface
func (r *restore) Status() model.JobMeta {
	r.stateM.Lock()
	jm := *r.job.Meta
	r.stateM.Unlock()
	return jm
}

// Cancel for jobber interface
func (r *restore) Cancel() {
	r.stateM.Lock()
	r.job.Meta.State = model.StateCanceling
	close(r.cancel)
	r.stateM.Unlock()
}

func (r *restore) SetState(state string) {
	r.stateM.Lock()
	r.job.Meta.State = state
	r.stateM.Unlock()
}

func (r *restore) GetState() string {
	r.stateM.Lock()
	s := r.job.Meta.State
	r.stateM.Unlock()
	return s
}

func (r *restore) addTotal(num int) {
	r.stateM.Lock()
	r.job.Meta.Total += num
	r.stateM.Unlock()
}

func (r *restore) addComplete(num int) {
	r.stateM.Lock()
	r.job.Meta.Complete += num
	r.stateM.Unlock()
}

func (r *restore) Run(finished chan<- string) {
	r.SetState(model.StateRunning)
	r.cancel = make(chan struct{})
	r.job.Meta.Start = time.Now()

	q := gowork.NewQueue(100, r.MaxWorkers)
	q.Start(r.MaxWorkers)

	completed := make(chan *model.JobFile, 20)
	errc := make(chan error)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		rclient, err := r.coordClient.Restore(context.Background(), &pb.RestoreRequest{Id: r.job.ID})
		if err != nil {
			errc <- err
			return
		}
		total := 0
		for {
			f, err := rclient.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				errc <- err
				return
			}

			total++
			if total > 10 {
				r.addTotal(total)
				total = 0
			}

			s := f.GetFile().GetSignature()
			m := f.GetFile().GetMeta()

			q.AddWork(work.Restore{
				File: files.File{
					Signature: files.Signature{Path: s.Path, Hash: s.Hash, Modifications: s.Modifications},
					Meta:      files.Meta{Mode: m.Mode, GID: int(m.Gid), UID: int(m.Uid)},
				},
				From:          *r.job.Definition.From,
				To:            r.job.Definition.To,
				Modifications: r.job.Definition.Modifications})
		}

		q.Finish()
	}()

	wg.Add(1)
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
				r.addComplete(processedFiles)
				processedFiles = 0
			}
		}
		r.addComplete(processedFiles)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fileClient, err := r.coordClient.File(context.Background())
		if err != nil {
			fmt.Printf("Error file client: %+v\n", err)
			errc <- err
			return
		}
		for file := range completed {
			fr := &pb.FileRequest{
				JobId: r.job.ID,
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
			fmt.Printf("%s : error: %s\n", r.job.ID, err)
		}
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
	finished <- r.job.ID
}
