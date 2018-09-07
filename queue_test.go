package queue

import (
	"fmt"
	"testing"
)

type mySimpler struct {
	i []int
}

func (m mySimpler) Len() int {
	return len(m.i)
}

func (m mySimpler) Pop(i int) interface{} {
	return m.i[i]
}

func TestCountChanQueue(t *testing.T) {
	s := mySimpler{i: make([]int, 100)}
	q := NewQueue(s)
	chCount := make(chan int)
	go func(ch chan int) {
		count := 0
		for {
			select {
			// case i := <-q.Pop():
			// 	fmt.Println("first", i)
			// 	count++
			// case i := <-q.Pop():
			// 	fmt.Println("second", i)
			// 	count++
			case i := <-q.Pop():
				fmt.Println("get", i)
				count++
			case <-q.Empty():
				fmt.Println("quit")
				ch <- count
			}
		}
	}(chCount)

	count := <-chCount

	if count != 100 {
		t.Error("it should return only 100 times but got", count)
	}
}
