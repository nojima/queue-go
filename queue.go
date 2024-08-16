package queue

import (
	"fmt"
	"iter"
	"math/bits"
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
	return q.length
}

// IsEmpty returns true if the queue is empty.
func (q *Queue[T]) IsEmpty() bool {
	return q.length == 0
}

// Push adds an element to the back of the queue.
func (q *Queue[T]) Push(x T) {
	if q.remainingCapacity() == 0 {
		q.reserve(len(q.buffer) + 1)
	}

	q.buffer[q.wrap(q.head+q.length)] = x
	q.length++
}

// PushMany adds multiple elements to the back of the queue.
// PushMany is more efficient than calling Push multiple times.
func (q *Queue[T]) PushMany(xs []T) {
	if q.remainingCapacity() < len(xs) {
		q.reserve(q.length + len(xs))
	}

	tail := q.wrap(q.head + q.length)
	n := copy(q.buffer[tail:], xs)
	copy(q.buffer, xs[n:])

	q.length += len(xs)
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

// remainingCapacity returns the number of elements that the buffer can still accommodate.
func (q *Queue[T]) remainingCapacity() int {
	return len(q.buffer) - q.length
}

// reserve ensures that the buffer has enough capacity to store requiredCapacity elements.
// Caller must guarantee that requiredCapacity > len(buffer).
func (q *Queue[T]) reserve(requiredCapacity int) {
	newCapacity := bitCeil(uint(requiredCapacity))

	capacity := len(q.buffer)
	if capacity == 0 {
		q.buffer = make([]T, newCapacity)
		return
	}

	newBuffer := make([]T, newCapacity)
	head := q.wrap(q.head)
	length := q.length
	n := copy(newBuffer, q.buffer[head:min(head+length, capacity)])
	copy(newBuffer[n:], q.buffer[:length-n])

	q.head = 0
	q.buffer = newBuffer
}

// bitCeil returns the minimum power of 2 that is greater than or equal to x.
// It returns 0 when x is 0.
func bitCeil(x uint) uint {
	return 1 << (bits.UintSize - bits.LeadingZeros(x-1))
}
