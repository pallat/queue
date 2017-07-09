package main

import (
	"fmt"

	"github.com/pallat/queue"
)

type lot struct {
	items  []int
	count  chan struct{}
	notify chan struct{}
}

func (l *lot) do(i <-chan int) {
	for {
		fmt.Println(l.items[<-i])
		l.count <- struct{}{}
		if len(l.count) == len(l.items) {
			l.notify <- struct{}{}
			return
		}
	}
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

	go l.do(q.Pop())
	go l.do(q.Pop())
	go l.do(q.Pop())
	go l.do(q.Pop())

	<-l.end()
	fmt.Println("total:", len(l.count))
}
