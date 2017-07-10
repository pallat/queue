package main

import (
	"fmt"

	"github.com/pallat/queue"
)

type lot struct {
	items  []int
	count  chan struct{}
	notify chan struct{}
	i      int
}

func (l *lot) do(i <-chan int) {
	for x := range i {
		fmt.Println(l.items[x])
		l.count <- struct{}{}
		l.i++
	}
}

func (l *lot) counting(x <-chan struct{}) {
	<-x
	for i := 0; i < len(l.items); i++ {
		<-l.count
	}
	l.notify <- struct{}{}
}

func (l *lot) end() <-chan struct{} {
	return l.notify
}

func main() {
	l := lot{
		count:  make(chan struct{}, 10),
		items:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		notify: make(chan struct{}),
	}

	q := queue.NewQueue(len(l.items))
	go l.counting(q.Empty())

	go l.do(q.Pop())
	go l.do(q.Pop())
	go l.do(q.Pop())
	go l.do(q.Pop())

	<-l.end()
	fmt.Println("total:", l.i)
}
