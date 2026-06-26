package queue

import (
	"sync"
)

// Thread safe bounded queue, which has a maximum size.
// When enqueueing, if full will reject and return false
type ThreadSafeQ struct {
	size int
	l    sync.Mutex
	v    []any
}

func NewTheadSafeQueue(size int) *ThreadSafeQ {
	return &ThreadSafeQ{size: size}
}

// Push onto the queue, return false if full
func (q *ThreadSafeQ) Push(e any) bool {
	q.l.Lock()
	defer q.l.Unlock()
	if len(q.v) == q.size {
		return false
	}
	q.v = append(q.v, e)
	return true
}

func (q *ThreadSafeQ) Pop() any {
	q.l.Lock()
	defer q.l.Unlock()
	if len(q.v) > 0 {
		e := q.v[len(q.v)-1]
		q.v = q.v[0 : len(q.v)-1]
		return e
	}
	return nil
}

func (q *ThreadSafeQ) Len() int {
	q.l.Lock()
	defer q.l.Unlock()
	return len(q.v)
}
