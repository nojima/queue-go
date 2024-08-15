package queue

// Queue is a FIFO queue backed by a circular buffer.
// The zero value for Queue is an empty queue ready to use.
type Queue[T any] struct {
	// Invariant: head <= tail (head == tail means the queue is empty)
	// Invariant: tail - head <= len(buffer)
	head int
	tail int

	// Invariant: len(buffer) is a power of 2 or zero
	buffer []T
}

// Len returns the number of elements in the queue.
func (q *Queue[T]) Len() int {
	return q.tail - q.head
}

// IsEmpty returns true if the queue is empty.
func (q *Queue[T]) IsEmpty() bool {
	return q.head == q.tail
}

// Push adds an element to the back of the queue.
func (q *Queue[T]) Push(x T) {
	q.growIfBufferIsFull()

	q.buffer[q.tailIndex()] = x
	q.tail++
}

// Pop removes and returns the element at the front of the queue.
func (q *Queue[T]) Pop() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}

	x := q.buffer[q.headIndex()]
	q.head++
	return x, true
}

// Peek returns the element at the front of the queue without removing it.
func (q *Queue[T]) Peek() (T, bool) {
	if q.Len() == 0 {
		var zero T
		return zero, false
	}

	return q.buffer[q.headIndex()], true
}

// headIndex returns the index of the head element in the buffer.
func (q *Queue[T]) headIndex() int {
	return q.head & (len(q.buffer) - 1)
}

// tailIndex returns the index where the next element would be added to the buffer.
func (q *Queue[T]) tailIndex() int {
	return q.tail & (len(q.buffer) - 1)
}

// growIfBufferIsFull double the buffer size if the buffer is full.
func (q *Queue[T]) growIfBufferIsFull() {
	length := q.Len()
	capacity := len(q.buffer)
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
	head := q.headIndex()
	n := copy(newBuffer, q.buffer[head:])
	copy(newBuffer[n:], q.buffer[:head])

	q.head = 0
	q.tail = length
	q.buffer = newBuffer
}
