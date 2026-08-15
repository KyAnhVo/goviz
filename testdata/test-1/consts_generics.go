package main

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

type Number interface {
	~int | ~int64 | ~float64
}

func Max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}
