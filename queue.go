package queue

// Simpler is the represent of your items, You just make you type with 2 methods
// Len() return the number of your items and Pop return each item by index
type Simpler interface {
	Pop(i int) interface{}
	Len() int
}

func NewQueue(s Simpler) *Queue {
	ch := make(chan int, s.Len())
	for i := 0; i < s.Len(); i++ {
		ch <- i
	}
	close(ch)

	q := &Queue{
		i:           ch,
		pop:         make(chan interface{}),
		emptyNotify: make(chan struct{}),
		s:           s,
	}
	go q.background()
	return q
}

type Queue struct {
	i           chan int
	pop         chan interface{}
	emptyNotify chan struct{}
	s           Simpler
}

func (q *Queue) background() {
	defer close(q.pop)
	defer close(q.emptyNotify)

	for i := range q.i {
		q.pop <- q.s.Pop(i)
	}
	q.emptyNotify <- struct{}{}
}

func (q *Queue) Pop() <-chan interface{} {
	return q.pop
}

func (q *Queue) Empty() <-chan struct{} {
	return q.emptyNotify
}
