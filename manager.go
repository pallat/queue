package queue

import (
	"context"
)

type Worker interface {
	Do(v interface{})
}

type Manager struct {
	count  chan struct{}
	notify chan struct{}
	i      int
	total  int
	w      Worker
	ctx    context.Context
	q      *Queue
	cherr  chan error
}

func NewManager(ctx context.Context, w Worker, s Simpler) *Manager {
	q := NewQueue(s)
	total := s.Len()
	m := &Manager{
		count:  make(chan struct{}, total),
		notify: make(chan struct{}),
		total:  total,
		i:      0,
		w:      w,
		ctx:    ctx,
		q:      q,
		cherr:  make(chan error),
	}
	go m.counting(q.Empty())
	return m
}

func (m *Manager) Do() {
	ch := make(chan interface{}, 1)
	for x := range m.q.Pop() {
		select {
		case ch <- x:
			m.w.Do(<-ch)
			m.count <- struct{}{}
			m.i++
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
	m.notify <- struct{}{}
}

func (m *Manager) End() <-chan struct{} {
	return m.notify
}
