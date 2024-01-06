package main

import (
	"fmt"
	"os"
)

func Die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func Assert(condition bool) {
	if !condition {
		panic("assertion failed")
	}
}
