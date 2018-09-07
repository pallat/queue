package queue

// Simpler is the represent of your items, You just make you type with 2 methods
// Len() return the number of your items and Pop return each item by index
type Simpler interface {
	Pop(i int) interface{}
	Len() int
}

func NewQueue(s Simpler) *Queue {
	q := &Queue{
		pop:         make(chan interface{}),
		emptyNotify: make(chan struct{}),
		s:           s,
	}
	go q.background()
	return q
}

type Queue struct {
	pop         chan interface{}
	emptyNotify chan struct{}
	s           Simpler
}

func (q *Queue) background() {
	defer close(q.emptyNotify)

	for i := 0; i < q.s.Len(); i++ {
		q.pop <- q.s.Pop(i)
	}
}

func (q *Queue) Pop() <-chan interface{} {
	return q.pop
}

func (q *Queue) Empty() <-chan struct{} {
	return q.emptyNotify
}
