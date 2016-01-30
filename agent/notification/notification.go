package notification

import (
	"time"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/util/log"
	"github.com/sethjback/gobl/util/try"
)

// Notification must be implemented on messages passed through the notifier
type Notification struct {
	Endpoint string
	Payload  []byte
}

// Queue handles queuing up messages to be sent to the coordinator
type Queue struct {
	Coordinator *config.Coordinator
	JobID       int
	Pending     []*Notification
	stop        chan struct{}
	In          chan *Notification
	next        chan *Notification
	Finished    chan bool
	waitTimer   *time.Timer
	keyManager  keys.Manager
}

// NewQueue creats and initilizes a new notifiaction queue
func NewQueue(c *config.Coordinator, jobID int, km keys.Manager) *Queue {
	return &Queue{
		Coordinator: c,
		JobID:       jobID,
		stop:        make(chan struct{}),
		In:          make(chan *Notification),
		next:        make(chan *Notification),
		Finished:    make(chan bool),
		keyManager:  km}
}

// QueueFromDisk loads a notification queue from disk
func QueueFromDisk(persistPath string) (*Queue, error) {
	return nil, nil
}

// SliceIQ is an infiniteQueue - Based on:
// github.com/kylelemons/iq
func (n *Queue) processQ() {

	log.Info("notification", "Process Q started")

recv:
	for {
		// Ensure that pending always has values so the select can
		// multiplex between the receiver and sender properly
		if len(n.Pending) == 0 {
			v, ok := <-n.In
			if !ok {
				// in is closed, flush values
				break
			}

			// We now have something to send
			n.Pending = append(n.Pending, v)
		}

		select {
		// Queue incoming values
		case v, ok := <-n.In:
			if !ok {
				// in is closed, flush values
				break recv
			}
			n.Pending = append(n.Pending, v)

			// Send queued values
		case n.next <- n.Pending[0]:
			n.Pending = n.Pending[1:]

		//stop closed, which means we need to exit without flushing
		case <-n.stop:
			log.Infof("notification", "Process Queue got stop, pending: %v", len(n.Pending))
			return
		}
	}

	// After in is closed, we may still have events to send
	log.Infof("notification", "Flushing queue. length: %v", len(n.Pending))
	for _, v := range n.Pending {
		select {
		case n.next <- v:
		case <-n.stop:
			//stop called...finish
			return
		}
	}

	//Lastly, we close the next channel to tell the notifier we are done
	close(n.next)
}

// Notify send the actual data to Cooridnator
// No guarantee that the data will be sent sequentually if errors occur when making the attempts
func (n *Queue) notify() {
	log.Info("notification", "Notify started")
	t := try.New(3)
	timeoutNum := 1
	var waitTime time.Duration
notify:
	for {
		select {
		case msg, ok := <-n.next:
			if !ok {
				//next was closed, we are done
				n.Finished <- true
				break notify
			}

			err := t.Do(func(attempt int) (bool, error) {
				sig, err := n.keyManager.Sign(string(msg.Payload))
				if err != nil {
					//log error
					log.Errorf("notification", "Unable to sign message: %v", err)
					return false, err
				}
				req := &httpapi.APIRequest{Address: n.Coordinator.Address, Body: msg.Payload, Signature: sig}

				_, err = req.POST(msg.Endpoint)
				if err != nil {
					log.Errorf("notification", "Unable to notify: %v", err)
					return false, err
				}

				return false, nil
			})

			if try.IsMaxRetries(err) {

				//requeue the message
				n.In <- msg

				//Communication issues, sleep then try again
				if timeoutNum < 10 {
					waitTime = time.Second * time.Duration(timeoutNum*timeoutNum)
				} else if timeoutNum < 15 {
					waitTime = time.Minute * time.Duration(timeoutNum)
				} else if timeoutNum < 21 {
					waitTime = time.Minute * time.Duration(timeoutNum*timeoutNum)
				} else {
					waitTime = time.Hour * time.Duration(timeoutNum)
					timeoutNum--
				}

				n.waitTimer = time.NewTimer(waitTime)
				<-n.waitTimer.C
				timeoutNum++
			}
		case <-n.stop:
			//Quit
			log.Info("notify", "Notify received stop call")
			break notify
		}
	}
}

// PersistQueue writes any pending notifications to disk
func (n *Queue) PersistQueue(persistDir string) {

}

// Run starts processing the q
func (n *Queue) Run() {
	log.Info("notification", "Starting")
	go n.processQ()
	go n.notify()
}

// Stop the queue from running
func (n *Queue) Stop() {
	close(n.stop)
	if n.waitTimer != nil {
		n.waitTimer.Stop()
	}
}

//Finish closes the in channel and starts a queue flush
func (n *Queue) Finish(f *Notification) {
	log.Info("notification", "Received finish")
	n.In <- f
	close(n.In)
}
