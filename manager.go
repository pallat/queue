package queue

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
}

func NewManager(q *Queue, w Worker, items ...interface{}) *Manager {
	total := len(items)
	m := &Manager{
		count:  make(chan struct{}, total),
		items:  items,
		notify: make(chan struct{}),
		total:  total,
		i:      0,
		w:      w,
	}
	go m.counting(q.Empty())
	return m
}

func (m *Manager) Do(i <-chan int) {
	for x := range i {
		m.w.Do(m.items[x])
		m.count <- struct{}{}
		m.i++
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
