package queue

import (
	"fmt"
	"iter"
)

// Queue is a FIFO queue backed by a circular buffer.
// The zero value for Queue is an empty queue ready to use.
// Queue is NOT safe for concurrent use.
type Queue[T any] struct {
	// Invariant: head <= tail (head == tail means the queue is empty)
	// Invariant: tail - head <= len(buffer)
	// Note that head and tail can be greater than len(buffer).
	head uint64 // a virtual index to read from
	tail uint64 // a virtual index to write to

	// Invariant: len(buffer) is a power of 2 or zero
	buffer []T
}

// Len returns the number of elements in the queue.
func (q *Queue[T]) Len() int {
	// Because builtin len() returns an int, q.Len() should return an int too.
	return int(q.tail - q.head)
}

// IsEmpty returns true if the queue is empty.
func (q *Queue[T]) IsEmpty() bool {
	return q.head == q.tail
}

// Push adds an element to the back of the queue.
func (q *Queue[T]) Push(x T) {
	q.growIfBufferIsFull()

	q.buffer[q.index(q.tail)] = x
	q.tail++
}

// Pop removes and returns the element at the front of the queue.
// If the queue is empty, Pop returns the zero value of T and false.
func (q *Queue[T]) Pop() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}

	x := q.buffer[q.index(q.head)]
	q.head++
	return x, true
}

// Peek returns the element at the front of the queue without removing it.
// If the queue is empty, Peek returns the zero value of T and false.
func (q *Queue[T]) Peek() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}

	return q.buffer[q.index(q.head)], true
}

// All returns an iterator over all elements in the queue.
// Do not modify the queue while iterating.
func (q *Queue[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := q.head; i < q.tail; i++ {
			if !yield(q.buffer[q.index(i)]) {
				break
			}
		}
	}
}

// Backward returns an iterator over all elements in the queue in reverse order (newest first).
// Do not modify the queue while iterating.
func (q *Queue[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := q.tail; i > q.head; i-- {
			if !yield(q.buffer[q.index(i-1)]) {
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
	return q.buffer[q.index(q.head+uint64(i))]
}

// index converts a virtual index into a buffer index.
func (q *Queue[T]) index(i uint64) uint64 {
	return i & (uint64(len(q.buffer)) - 1)
}

// growIfBufferIsFull double the buffer size if the buffer is full.
func (q *Queue[T]) growIfBufferIsFull() {
	length := q.tail - q.head
	capacity := uint64(len(q.buffer))
	if length < capacity {
		return
	}
	if capacity == 0 {
		q.head = 0
		q.tail = 0
		q.buffer = make([]T, 1)
		return
	}

	newBuffer := make([]T, capacity*2)
	head := q.index(q.head)
	n := copy(newBuffer, q.buffer[head:])
	copy(newBuffer[n:], q.buffer[:head])

	q.head = 0
	q.tail = length
	q.buffer = newBuffer
}
