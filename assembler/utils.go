package main

import (
	"fmt"
	"os"
)

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

func Die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
