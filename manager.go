package queue

import (
	"context"
)

type Worker interface {
	Do(v interface{}) error
}

type Manager struct {
	count  chan struct{}
	notify chan struct{}
	total  int
	w      Worker
	ctx    context.Context
	q      *Queue
	chv    chan interface{}
}

func NewManager(ctx context.Context, w Worker, s Simpler) *Manager {
	q := NewQueue(s)
	total := s.Len()
	m := &Manager{
		count:  make(chan struct{}, total),
		notify: make(chan struct{}),
		total:  total,
		w:      w,
		ctx:    ctx,
		q:      q,
		chv:    make(chan interface{}, total),
	}
	go m.counting(q.Empty())
	return m
}

func (m *Manager) Execute(n int) {
	for i := 0; i < n; i++ {
		go m.Do()
	}

	<-m.End()
}

func (m *Manager) Do() {
	ch := make(chan interface{}, 1)
	for x := range m.q.Pop() {
		select {
		case ch <- x:
			m.chv <- m.w.Do(<-ch)
			m.count <- struct{}{}
		case <-m.ctx.Done():
			m.notify <- struct{}{}
		}
	}
}

func (m *Manager) counting(x <-chan struct{}) {
	<-x
	for i := 0; i < m.total; i++ {
		<-m.count
	}
	close(m.chv)
	m.notify <- struct{}{}
}

func (m *Manager) End() <-chan struct{} {
	return m.notify
}

func (m *Manager) Response() <-chan interface{} {
	return m.chv
}
