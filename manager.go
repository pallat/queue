package queue

import (
	"context"
	"fmt"
)

type Worker interface {
	Do(v interface{})
}

type Manager struct {
	items  []interface{}
	count  chan struct{}
	notify chan struct{}
	i      int
	total  int
	w      Worker
	ctx    context.Context
}

func NewManager(ctx context.Context, q *Queue, w Worker, items ...interface{}) *Manager {
	total := len(items)
	m := &Manager{
		count:  make(chan struct{}, total),
		items:  items,
		notify: make(chan struct{}),
		total:  total,
		i:      0,
		w:      w,
		ctx:    ctx,
	}
	go m.counting(q.Empty())
	return m
}

func (m *Manager) Do(i <-chan int) {
	ch := make(chan int, 1)
	for x := range i {
		select {
		case ch <- x:
			m.w.Do(m.items[<-ch])
			m.count <- struct{}{}
			m.i++
		case <-m.ctx.Done():
			fmt.Println(m.ctx.Err)
			m.notify <- struct{}{}
		}
	}
}

func (m *Manager) counting(x <-chan struct{}) {
	<-x
	for i := 0; i < len(m.items); i++ {
		<-m.count
	}
	m.notify <- struct{}{}
}

func (m *Manager) End() <-chan struct{} {
	return m.notify
}
