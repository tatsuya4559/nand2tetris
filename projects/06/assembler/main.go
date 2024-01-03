package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func assertGivenFileIsAssembly(filename string) {
	if filepath.Ext(filename) != ".asm" {
		fmt.Fprintf(os.Stderr, "File must be .asm file\n")
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	filename := flag.Arg(0)

	assertGivenFileIsAssembly(filename)

	asmFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer asmFile.Close()

	parser := NewParser(asmFile)
	for parser.HasMoreCommand() {
		parser.Advance()
		fmt.Println(parser.CurrentCommand())
	}
}
