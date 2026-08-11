package lexer

type Queue[T any] struct {
	buf  []T
	head int
	size int
}

func NewQueue[T any](capacity int) *Queue[T] {
	if capacity <= 0 {
		panic("capacity must be positive")
	}

	return &Queue[T]{
		buf: make([]T, capacity),
	}
}

func (q *Queue[T]) Len() int {
	return q.size
}

func (q *Queue[T]) Cap() int {
	return len(q.buf)
}

func (q *Queue[T]) Empty() bool {
	return q.size == 0
}

func (q *Queue[T]) Push(x T) {
	if q.size == len(q.buf) {
		q.grow()
	}

	tail := (q.head + q.size) % len(q.buf)
	q.buf[tail] = x
	q.size++
}

func (q *Queue[T]) Pop() T {
	if q.size == 0 {
		panic("pop from empty queue")
	}

	x := q.buf[q.head]

	// Clear the slot so references can be garbage-collected.
	var zero T
	q.buf[q.head] = zero

	q.head = (q.head + 1) % len(q.buf)
	q.size--

	return x
}

// At returns the i-th element in queue order.
// At(0) is the front, At(Len()-1) is the back.
func (q *Queue[T]) At(i int) T {
	if i < 0 || i >= q.size {
		panic("queue index out of range")
	}

	return q.buf[(q.head+i)%len(q.buf)]
}

// AtPtr returns a pointer to the i-th element.
// The pointer becomes invalid if the queue grows or the element
// is removed.
func (q *Queue[T]) AtPtr(i int) *T {
	if i < 0 || i >= q.size {
		panic("queue index out of range")
	}

	return &q.buf[(q.head+i)%len(q.buf)]
}

func (q *Queue[T]) grow() {
	newBuf := make([]T, len(q.buf)*2)

	for i := 0; i < q.size; i++ {
		newBuf[i] = q.buf[(q.head+i)%len(q.buf)]
	}

	q.buf = newBuf
	q.head = 0
}
