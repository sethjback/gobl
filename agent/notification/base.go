package notification

import (
	"crypto/rsa"
	"net/http"
	"sync"
	"time"

	"github.com/eapache/queue"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/util/log"
	"github.com/sethjback/gowork"
)

// baseNotifier implements Notifier
type baseNotifier struct {
	signingKey *rsa.PrivateKey
	hclient    *http.Client
	config     *Config
	pending    *queue.Queue
	retry      *queue.Queue
	waiter     *sync.WaitGroup
	in         chan *Message
	rin        chan *Message
	send       chan *Message
	result     chan *Result
	stop       bool
	stopLock   *sync.Mutex
}

func newBase(config *Config, key *rsa.PrivateKey) *baseNotifier {
	if key == nil {
		panic("Must provide rsa private key")
	}

	if config == nil {
		config = &Config{MaxWorkers: 3, MaxDepth: 20}
	}

	n := &baseNotifier{
		config:     config,
		signingKey: key,
		hclient:    &http.Client{},
		pending:    queue.New(),
		retry:      queue.New(),
		waiter:     &sync.WaitGroup{},
		in:         make(chan *Message),
		rin:        make(chan *Message),
		send:       make(chan *Message),
		result:     make(chan *Result),
		stopLock:   &sync.Mutex{},
	}

	n.waiter.Add(3)
	return n
}

func (n *baseNotifier) Start() {
	go n.manageQs()
	go n.manageResults()
	go n.manageSender()
}

// true of the notifier is in the stopped state.
func (n *baseNotifier) Stopped() bool {
	var s bool
	n.stopLock.Lock()
	s = n.stop
	n.stopLock.Unlock()
	return s
}

// Stop the notifier. This will allow pending notifications to finish then flush any queued notifications to disk
func (n *baseNotifier) Stop() {
	n.stopLock.Lock()
	n.stop = true
	n.stopLock.Unlock()

	//todo, wait on the waiter then flush pending and retry queue to disk
	n.waiter.Wait()
}

// Send a message.
func (n *baseNotifier) Send(note Notification) {
	if !n.Stopped() {
		n.in <- &Message{retry: 0, note: note}
	}
}

// infinite channel queue
// inspired by https://godoc.org/github.com/eapache/channels#InfiniteChannel
func (n *baseNotifier) manageQs() {
	log.Info("notifier", "managing queues")

	// signal we are done
	defer n.waiter.Done()

	// every 2 minutes check the retry queue
	t := time.NewTicker(2 * time.Minute)

	var next *Message
	var in, retry, send chan *Message

	in = n.in
	retry = n.rin

	for in != nil || retry != nil || send != nil {
		select {
		case m, open := <-in:
			if open {
				n.pending.Add(m)
			} else {
				in = nil
			}

		case m, open := <-retry:
			if open {
				n.retry.Add(m)
			} else {
				retry = nil
			}

		case send <- next:
			n.pending.Remove()

		case <-t.C:
			//iterate over retry queue, add to pending if appropriate
			clength := n.retry.Length()
			for i := 0; i < clength; i++ {
				m := n.retry.Remove().(*Message)
				if m.retry == 0 {
					n.pending.Add(m)
				} else {
					m.retry--
					n.retry.Add(m)
				}
			}

		}

		if n.Stopped() {
			// closing send will start the shut down process:
			// it will cause the sending queue to stop and drain
			log.Debug("notifier", "We are stopped")
			if send != nil {
				close(send)
			}
			send = nil
			if in != nil {
				close(in)
			}
			in = nil
			next = nil
		} else {
			if n.pending.Length() > 0 {
				send = n.send
				next = n.pending.Peek().(*Message)
			} else {
				send = nil
				next = nil
			}
		}
	}

}

func (n *baseNotifier) manageResults() {
	log.Info("notifier", "managing retry")

	// signal we are done
	defer n.waiter.Done()

	for r := range n.result {
		switch r.state {
		case Retry:
			if r.message.retry < 5 {
				r.message.retry++
			}
			n.rin <- r.message
		case Fail:
			log.Errorf("notifier", "Unable send message to: %s || %s", r.message.note.Host()+"/"+r.message.note.Path(), string(r.message.note.Body()))
		default:
			//move on
		}
	}

	log.Debug("notifier", "result close, close result in")

	close(n.rin)
}

func (n *baseNotifier) manageSender() {
	log.Info("notifier", "managing sender")

	q := gowork.NewQueue(n.config.MaxDepth, n.config.MaxDepth)
	q.Start(n.config.MaxWorkers)

	defer n.waiter.Done()

	go func() {
		for m := range n.send {
			q.AddWork(&Sender{
				client:  n.hclient,
				signer:  keys.NewSigner(n.signingKey),
				message: m,
			})
		}
		log.Debug("notifier", "Send Closed, finish")
		q.Finish()
	}()

	done := make(chan struct{})

	go func() {
		for r := range q.Results() {
			n.result <- r.(*Result)
		}
		log.Debug("notifier", "results closed. close done")
		close(done)
	}()

	// wait for the queue to finish
	<-done

	log.Debug("notifier", "done done, close results")
	close(n.result)
}
