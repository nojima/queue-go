package queue_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/nojima/queue-go"
)

func ExampleQueue() {
	var q queue.Queue[int]
	q.Push(3)
	q.Push(1)
	q.Push(4)

	for !q.IsEmpty() {
		x, ok := q.Pop()
		if !ok {
			panic("queue should not be empty here")
		}
		fmt.Println(x)
	}
	// Output:
	// 3
	// 1
	// 4
}

func ExampleQueue_All() {
	var q queue.Queue[int]
	q.Push(3)
	q.Push(1)
	q.Push(4)

	for x := range q.All() {
		fmt.Println(x)
	}
	// Output:
	// 3
	// 1
	// 4
}

func TestRandomized(t *testing.T) {
	var q queue.Queue[int]
	var v []int

	for i := 0; i < 10000; i++ {
		r := rand.Uint32()
		if r%2 == 0 {
			q.Push(i)
			v = append(v, i)
		} else {
			x, ok := q.Pop()

			var expectedX int
			var expectedOK bool
			if len(v) == 0 {
				expectedX = 0
				expectedOK = false
			} else {
				expectedX = v[0]
				expectedOK = true
				v = v[1:]
			}

			if x != expectedX || ok != expectedOK {
				t.Errorf("Pop() = %v, %v; want %v, %v", x, ok, expectedX, expectedOK)
			}
		}

		if q.Len() != len(v) {
			t.Errorf("Len() = %v; want %v", q.Len(), len(v))
		}

		x, ok := q.Peek()
		var expectedX int
		var expectedOK bool
		if len(v) == 0 {
			expectedX = 0
			expectedOK = false
		} else {
			expectedX = v[0]
			expectedOK = true
		}
		if x != expectedX || ok != expectedOK {
			t.Errorf("Peek() = %v, %v; want %v, %v", x, ok, expectedX, expectedOK)
		}
	}
}
