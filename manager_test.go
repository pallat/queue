package queue

import (
	"fmt"
	"testing"
)

func TestQueuManager(t *testing.T) {
	items := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var w work = 0
	q := NewQueue(len(items))
	m := NewManager(q, &w, items...)

	go m.Do(q.Pop())
	go m.Do(q.Pop())
	go m.Do(q.Pop())
	go m.Do(q.Pop())

	<-m.End()

	if w != 10 {
		t.Error("not finish")
	}
}

type work int

func (w *work) Do(v interface{}) {
	fmt.Println(w, v)
	*w++
}
