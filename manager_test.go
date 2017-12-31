package queue

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type work int

func (w *work) Do(v interface{}) {
	fmt.Printf("-->%#v\n%#v\n", w, v)
	*w++
}

func TestQueuManager(t *testing.T) {
	s := mySimpler{i: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}

	var w work
	ctx := context.Background()
	m := NewManager(ctx, &w, s)

	go m.Do()
	go m.Do()
	go m.Do()
	go m.Do()
	go m.Do()
	go m.Do()

	<-m.End()

	if w != 10 {
		t.Error("not finish", w)
	}
}

// func TestQueuManagerWith1To50Fields(t *testing.T) {
// 	items := []interface{}{0, 1, 2}
// 	s := mySimpler{i: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}

// 	for i := 3; i < 50; i++ {
// 		items = append(items, i)
// 		total := len(items)
// 		var w work
// 		ctx := context.Background()
// 		m := NewManager(ctx, &w, items...)

// 		go m.Do()
// 		go m.Do()
// 		go m.Do()
// 		go m.Do()

// 		<-m.End()

// 		if w != work(total) {
// 			t.Error("not finish", w)
// 			return
// 		}
// 	}
// }

// func TestQueuManagerWithIntItems(t *testing.T) {
// 	var w work
// 	ctx := context.Background()
// 	m := NewManager(ctx, &w, 1, 2, 3, 4, 5)

// 	go m.Do()
// 	go m.Do()
// 	go m.Do()
// 	go m.Do()
// 	go m.Do()
// 	go m.Do()

// 	<-m.End()

// 	if w != 5 {
// 		t.Error("not finish", w)
// 	}
// }

type slow int

func (w *slow) Do(v interface{}) {
	time.Sleep(500 * time.Millisecond)
	*w++
	fmt.Println(w, v)
}

func TestQueuManagerTimeout(t *testing.T) {
	// items := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s := mySimpler{i: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}

	var w slow
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	m := NewManager(ctx, &w, s)

	go m.Do()
	go m.Do()
	go m.Do()
	go m.Do()

	<-m.End()

	if w == 10 {
		t.Error("not timeout", w)
	}
}
