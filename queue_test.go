package queue

import (
	"fmt"
	"testing"
)

func TestCountChanQueue(t *testing.T) {
	q := NewQueue(10)
	chCount := make(chan int)
	go func(ch chan int) {
		count := 0
		for {
			select {
			case i := <-q.Pop():
				fmt.Println("first", i)
				count++
			case i := <-q.Pop():
				fmt.Println("second", i)
				count++
			case i := <-q.Pop():
				fmt.Println("third", i)
				count++
			case <-q.Empty():
				fmt.Println("quit")
				ch <- count
			}
		}
	}(chCount)

	count := <-chCount

	if count != 10 {
		t.Error("it should return only 10 times but got", count)
	}
}
