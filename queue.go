package queue

import (
	"fmt"
	"iter"
)

// Queue is a FIFO queue backed by a circular buffer.
// The zero value for Queue is an empty queue ready to use.
// Queue is NOT safe for concurrent use.
type Queue[T any] struct {
	// The index of the first element in the queue.
	// Invariant:
	//   0 <= head < len(buffer)  (len(buffer) != 0)
	//   head == 0                (len(buffer) == 0)
	head int

	// The number of elements in the queue.
	// Invariant: length <= len(buffer)
	length int

	// The circular buffer to store elements.
	// Invariant: len(buffer) is a power of 2 or zero
	buffer []T
}

// Len returns the number of elements in the queue.
func (q *Queue[T]) Len() int {
	// Because builtin len() returns an int, q.Len() should return an int too.
	return q.length
}

// IsEmpty returns true if the queue is empty.
func (q *Queue[T]) IsEmpty() bool {
	return q.length == 0
}

// Push adds an element to the back of the queue.
func (q *Queue[T]) Push(x T) {
	if q.isBufferFull() {
		q.grow()
	}

	q.buffer[q.wrap(q.head+q.length)] = x
	q.length++
}

// Pop removes and returns the element at the front of the queue.
// If the queue is empty, Pop returns the zero value of T and false.
func (q *Queue[T]) Pop() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}

	x := q.buffer[q.head]
	q.head = q.wrap(q.head + 1)
	q.length--
	return x, true
}

// Peek returns the element at the front of the queue without removing it.
// If the queue is empty, Peek returns the zero value of T and false.
func (q *Queue[T]) Peek() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}

	return q.buffer[q.head], true
}

// All returns an iterator over all elements in the queue.
// Do not modify the queue while iterating.
func (q *Queue[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		head := q.head
		for i := range q.length {
			if !yield(q.buffer[q.wrap(head+i)]) {
				break
			}
		}
	}
}

// Backward returns an iterator over all elements in the queue in reverse order (newest first).
// Do not modify the queue while iterating.
func (q *Queue[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		last := q.head + q.length - 1
		for i := range q.length {
			if !yield(q.buffer[q.wrap(last-i)]) {
				break
			}
		}
	}
}

// At returns the element at the specified index.
// If the index is out of range, it panics.
func (q *Queue[T]) At(i int) T {
	if i < 0 || i >= q.Len() {
		panic(fmt.Sprintf("queue: index out of range: i=%d, len=%d", i, q.Len()))
	}
	return q.buffer[q.wrap(q.head+i)]
}

// wrap converts an index to the corresponding index in the buffer.
func (q *Queue[T]) wrap(i int) int {
	return i & (len(q.buffer) - 1)
}

// isBufferFull returns true if the buffer is full.
func (q *Queue[T]) isBufferFull() bool {
	return len(q.buffer) == q.length
}

// grow doubles the buffer size.
// This method is called when the buffer is full.
func (q *Queue[T]) grow() {
	capacity := len(q.buffer)
	if capacity == 0 {
		q.buffer = make([]T, 1)
		return
	}

	newBuffer := make([]T, capacity*2)
	head := q.wrap(q.head)
	n := copy(newBuffer, q.buffer[head:])
	copy(newBuffer[n:], q.buffer[:head])

	q.head = 0
	q.buffer = newBuffer
}
