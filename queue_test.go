package queue_test

import (
	"fmt"
	"math/rand"
	"slices"
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

func TestQueue_Backward(t *testing.T) {
	testCases := []struct {
		title    string
		elements []int
	}{
		{
			title:    "empty",
			elements: []int{},
		},
		{
			title:    "one element",
			elements: []int{1},
		},
		{
			title:    "multiple elements",
			elements: []int{3, 1, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// Setup
			var q queue.Queue[int]
			for _, x := range tc.elements {
				q.Push(x)
			}

			// Exercise
			var actual []int
			for x := range q.Backward() {
				actual = append(actual, x)
			}

			// Verify
			expected := slices.Clone(tc.elements)
			slices.Reverse(expected)
			if !slices.Equal(actual, expected) {
				t.Errorf("actual: %v; want: %v", actual, expected)
			}
		})
	}
}

func TestQueue_At(t *testing.T) {
	var q queue.Queue[int]
	for _, x := range []int{3, 1, 4, 1, 5, 9, 2} {
		q.Push(x)
	}
	q.Pop()
	q.Pop()

	for i, expected := range []int{4, 1, 5, 9, 2} {
		actual := q.At(i)
		if actual != expected {
			t.Errorf("At(%v) = %v; want %v", i, actual, expected)
		}
	}
}

func TestRandomized(t *testing.T) {
	for k := 0; k < 1000; k++ {
		var q queue.Queue[int]
		var v []int

		for i := 0; i < 1000; i++ {
			r := rand.Uint32()
			switch r % 3 {
			case 0:
				q.Push(i)
				v = append(v, i)
			case 1:
				var xs []int
				for range rand.Intn(10) {
					xs = append(xs, rand.Int())
				}
				q.PushMany(xs)
				v = append(v, xs...)
			case 2:
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

			if !q.IsEmpty() {
				x := q.At(q.Len() - 1)
				expectedX := v[len(v)-1]
				if x != expectedX {
					t.Errorf("At(%v) = %v; want %v", q.Len()-1, x, expectedX)
				}
			}
		}
	}
}

func BenchmarkPushPop(b *testing.B) {
	var q queue.Queue[int]
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
	var expected int
	for !q.IsEmpty() {
		x, ok := q.Pop()
		if !ok {
			b.Fatal("queue should not be empty here")
		}
		if x != expected {
			b.Errorf("Pop() = %v; want %v", x, expected)
		}
		expected++
	}
}
