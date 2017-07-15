package main

import (
	"fmt"

	"github.com/pallat/queue"
)

func main() {
	items := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	q := queue.NewQueue(len(items))
	m := queue.NewManager(q, work, items...)

	go m.Do(q.Pop())
	go m.Do(q.Pop())
	go m.Do(q.Pop())
	go m.Do(q.Pop())

	<-m.End()
}

func work(v interface{}) {
	fmt.Println(v)
}
