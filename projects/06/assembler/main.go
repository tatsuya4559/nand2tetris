package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func assertGivenFileIsAssembly(filename string) {
	if filepath.Ext(filename) != ".asm" {
		Die("File must be .asm file")
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

	// 1pass
	parser := NewParser(asmFile)
	symbolTable := NewSymbolTable()
	symbolTable.LoadLabelAddress(parser)

	hackFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".hack"
	hackFile, err := os.Create(hackFilename)
	if err != nil {
		Die("failed to create hack file: %v", err)
	}
	defer hackFile.Close()

	// 2pass
	asmFile.Seek(0, io.SeekStart)
	parser = NewParser(asmFile)
	for parser.Parse() {
		var bin uint16
		var err error
		switch cmd := parser.CurrentCommand().(type) {
		case *ACommand:
			if cmd.SymbolIsDigit {
				bin, err = ConvertACommand(cmd)
			} else if addr, ok := symbolTable.GetAddress(cmd.Symbol); ok {
				bin = addr
			} else {
				bin = symbolTable.AddAutoEntry(cmd.Symbol)
			}
		case *CCommand:
			bin, err = ConvertCCommand(cmd)
		default:
			continue
		}
		if err != nil {
			Die("failed to convert command to binary: %v", err)
		}
		fmt.Fprintf(hackFile, "%016b\n", bin)
	}
}
