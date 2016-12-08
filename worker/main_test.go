package worker

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testWorker struct {
	id    int
	c     *counter
	sleep time.Duration
}

type counter struct {
	mtx   *sync.Mutex
	count int
}

func (c *counter) Add(a int) {
	c.mtx.Lock()
	c.count += a
	c.mtx.Unlock()
}

func (c *counter) GetCount() int {
	c.mtx.Lock()
	i := c.count
	c.mtx.Unlock()
	return i
}

func (t *testWorker) Do() interface{} {
	time.Sleep(t.sleep)
	t.c.Add(1)
	return nil
}

func TestQueue(t *testing.T) {
	assert := assert.New(t)
	c := &counter{mtx: &sync.Mutex{}}
	q := NewQueue(100, 4)
	q.Start(4)

	for i := 0; i < 40; i++ {
		q.AddWork(&testWorker{sleep: 100 * time.Millisecond, id: i, c: c})
	}

	qAddTime := time.Now().Unix()

	q.Finish()
	for r := range q.ResultChan() {
		assert.Nil(r)
	}
	qFinishTime := time.Now().Unix()

	assert.NotEqual(qAddTime, qFinishTime)
	assert.Equal(40, c.GetCount())

	c = &counter{mtx: &sync.Mutex{}}
	q = NewQueue(100, 4)
	q.Start(4)
	for i := 0; i < 40; i++ {
		q.AddWork(&testWorker{sleep: 100 * time.Millisecond, id: i, c: c})
	}

	q.Finish()

	i := 0
	for r := range q.ResultChan() {
		assert.Nil(r)
		if i == 15 {
			q.Abort()
		}
		i++
	}

	assert.NotEqual(40, c.GetCount())
}
