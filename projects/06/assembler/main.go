package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	hackFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".hack"
	hackFile, err := os.Create(hackFilename)
	if err != nil {
		panic(err)
	}
	defer hackFile.Close()

	parser := NewParser(asmFile)
	for parser.HasMoreCommand() {
		scanned, err := parser.Advance()
		if err != nil {
			log.Fatal(err)
		}
		if scanned {
			bin, err := CommandToBinaryCode(parser.CurrentCommand())
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(hackFile, "%016b\n", bin)
		}
	}
}
