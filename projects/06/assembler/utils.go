package main

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](elem ...T) Set[T] {
	set := make(Set[T])
	for _, e := range elem {
		set[e] = struct{}{}
	}
	return set
}

func (s Set[T]) Contains(elem T) bool {
	_, ok := s[elem]
	return ok
}

func Assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}
