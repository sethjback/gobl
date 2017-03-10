package job

import (
	"sync"

	"github.com/sethjback/gobl/agent/notification"
)

// tn implements notification.Notifier for use testing backup/restore jobs
type tn struct {
	m       *sync.Mutex
	sent    []notification.Notification
	started bool
}

func newTestNotifier() *tn {
	return &tn{
		m:       &sync.Mutex{},
		started: false,
		sent:    []notification.Notification{}}
}

func (t *tn) Start() {
	t.started = true
}
func (t *tn) Stop() {
	t.started = false
}
func (t *tn) Stopped() bool {
	return t.started
}
func (t *tn) Send(note notification.Notification) {
	t.m.Lock()
	t.sent = append(t.sent, note)
	t.m.Unlock()
}
