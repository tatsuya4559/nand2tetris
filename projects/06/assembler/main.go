package main

import (
	"flag"
	"fmt"
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
		Die("failed to open asm file: %v", err)
	}
	defer asmFile.Close()

	hackFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".hack"
	hackFile, err := os.Create(hackFilename)
	if err != nil {
		Die("failed to create hack file: %v", err)
	}
	defer hackFile.Close()

	parser := NewParser(asmFile)
	for parser.Parse() {
		var bin uint16
		var err error
		switch cmd := parser.CurrentCommand().(type) {
		case *ACommand:
			bin, err = ConvertACommand(cmd)
		case *CCommand:
			bin, err = ConvertCCommand(cmd)
		default:
			Die("Command %v does not have binary representation", cmd)
		}
		if err != nil {
			Die("failed to convert command to binary: %v", err)
		}
		fmt.Fprintf(hackFile, "%016b\n", bin)
	}
}
