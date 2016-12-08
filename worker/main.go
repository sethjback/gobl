package worker

// Work is the interface that must be implemented to run on the work queue
type Work interface {
	// Do any appropriate work, returning an interface to be passed back on the result channel
	Do() interface{}
}
