package util

import "errors"

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](lst []T) Set[T] {
	set := make(Set[T])
	for _, elem := range lst {
		set[elem] = struct{}{}
	}
	return set
}

func (s Set[T]) Contains(val T) bool {
	_, ok := s[val]
	return ok
}

type Queue[T any] struct {
	buf  []T
	cap  int
	len  int
	head int
	tail int
}

func NewQueue[T any](cap int) *Queue[T] {
	return &Queue[T]{
		buf:  make([]T, cap),
		cap:  cap,
		len:  0,
		head: 0,
	}
}

func (q *Queue[T]) At(loc int) (T, error) {
	var val T
	if q.len <= loc {
		return val, errors.New("loc > q size")
	}
	return q.buf[(q.head+loc)%q.cap], nil
}

func (q *Queue[T]) Enqueue(val T) error {
	if q.len == q.cap {
		return errors.New("exceed length")
	}
	q.buf[(q.head+q.len)%q.cap] = val
	q.len += 1
	return nil
}

func (q *Queue[T]) Dequeue() (T, error) {
	var val T
	if q.len == 0 {
		return val, errors.New("empty buffer")
	}
	val = q.buf[q.head]
	q.head = (q.head + 1) % q.cap
	q.len -= 1
	return val, nil
}

func (q *Queue[T]) Push(val T) error {
	if q.len == q.cap {
		return errors.New("exceed length")
	}
	q.head = (q.head + q.cap - 1) % q.cap
	q.buf[q.head] = val
	q.len += 1
	return nil
}

func (q *Queue[T]) Pop() (T, error) {
	var val T
	if q.len == 0 {
		return val, errors.New("empty buffer")
	}
	val = q.buf[q.head]
	q.head = (q.head + 1) % q.cap
	q.len -= 1
	return val, nil
}
