package queue

func NewQueue(total int) *Queue {
	ch := make(chan int, total)
	for i := 0; i < total; i++ {
		ch <- i
	}

	q := &Queue{
		i:          ch,
		pop:        make(chan int),
		fullNotify: make(chan struct{}),
	}
	go q.background()
	return q
}

type Queue struct {
	i          chan int
	pop        chan int
	fullNotify chan struct{}
}

func (q *Queue) background() {
	for {
		q.pop <- <-q.i
		if len(q.i) == 0 {
			q.fullNotify <- struct{}{}
		}
	}
}

func (q *Queue) Pop() <-chan int {
	return q.pop
}

func (q *Queue) Full() <-chan struct{} {
	return q.fullNotify
}

func (q *Queue) Close() {
	close(q.fullNotify)
}
