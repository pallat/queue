package queue

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

type work int

func (w *work) Do(v interface{}) error {
	fmt.Printf("-->%#v\n%#v\n", w, v)
	*w++
	return nil
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

type slow int

func (w *slow) Do(v interface{}) error {
	time.Sleep(500 * time.Millisecond)
	*w++
	fmt.Println(w, v)
	return nil
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

type workPanic int

func (w *workPanic) Do(v interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	panic(errors.New("fail"))
	return nil
}

func TestQueuManagerWhenWorkerHasPanic(t *testing.T) {
	s := mySimpler{i: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}

	var w workPanic
	ctx := context.Background()
	m := NewManager(ctx, &w, s)

	go m.Do()
	go m.Do()

	<-m.End()

	if w != 0 {
		t.Error("not finish", w)
	}
}

type workError int

func (w *workError) Do(v interface{}) error {
	if (v.(int) % 2) == 0 {
		return errors.New("it's error")
	}
	return nil
}

func TestQueuManagerWhenWorkerHasError(t *testing.T) {
	s := mySimpler{i: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}

	var w workError
	ctx := context.Background()
	m := NewManager(ctx, &w, s)

	go m.Do()
	go m.Do()

	<-m.End()

	if w != 0 {
		t.Error("not finish", w)
	}

	for err := range m.Response() {
		fmt.Println("---->", err)
	}
}

func TestQueuManagerAssignGoRoutineNumber(t *testing.T) {
	s := mySimpler{i: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}

	var w work
	ctx := context.Background()
	m := NewManager(ctx, &w, s)

	m.Execute(12)

	if w != 10 {
		t.Error("not finish", w)
	}
}
