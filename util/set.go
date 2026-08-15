package util

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
