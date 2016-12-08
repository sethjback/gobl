package worker

import "sync"

// Queue is the collection point for pending work.
// It manages running incoming work on go routines, stops them
// if aborted, and cleaning up after they are finished
type Queue struct {
	workQueue  chan Work
	resultChan chan interface{}
	abortChan  chan struct{}
}

// NewQueue returns a new work queue
// workDepth represents how large the queue is for incoming work
// resultDepth represents how large the result queue is. Don't set the
// resultDepth too high - probably appropriate to set at the same number of the
// workers you will have running.
func NewQueue(workDepth, resultDepth int) *Queue {
	return &Queue{
		workQueue:  make(chan Work, workDepth),
		resultChan: make(chan interface{}, resultDepth),
		abortChan:  make(chan struct{}),
	}
}

// ResultChan resturns the channel all workers will be sending their results over.
// this must be ranged over - after the initial depth fills up workers will block
// until they are able to send over this channel
func (q *Queue) ResultChan() <-chan interface{} {
	return q.resultChan
}

// Push work onto the queue
// Must NOT call this after you have called Finish()
func (q *Queue) AddWork(w Work) {
	q.workQueue <- w
}

// Abort Will cause all workers to exit after finishing their current task
// This aborts processing the queue: any unfinished work will be ignored.
func (q *Queue) Abort() {
	close(q.abortChan)
}

// Finish() signals the queue that it is no longer needed. This MUST be called
// when you are through with the queue otherwise the worker threads will continue
// to listen for work. It is safe to call before all work is complete: workers
// will process any backlog in the work queue before exiting
func (q *Queue) Finish() {
	close(q.workQueue)
}

// Start spins up workerCount workers
// workers will continue to wait for work until Finish() is called.
// Once all workes have exited, the resultChan will be closed indicating that
// that all workers have finished gracefully
func (q *Queue) Start(workerCount int) {
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			doWork(q.workQueue, q.abortChan, q.resultChan)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(q.resultChan)
	}()
}

// doWork is the actual function that handles the work
func doWork(q <-chan Work, abort <-chan struct{}, resultChan chan<- interface{}) {
	for work := range q {
		select {
		case resultChan <- work.Do():
		case <-abort:
			return
		}
	}
}
