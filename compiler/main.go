package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		Die("Usage: %s [FILE | DIR]", os.Args[0])
	}
	path := flag.Arg(0)

	fileInfo, err := os.Stat(path)
	if err != nil {
		Die("cannot stat %s: %v", path, err)
	}
	if fileInfo.IsDir() {
		jackFiles, err := filepath.Glob(filepath.Join(path, "*.jack"))
		if err != nil {
			Die("jack files not found: %v", err)
		}
		for _, jackFile := range jackFiles {
			compileFile(jackFile)
		}
	} else {
		compileFile(path)
	}
}

func outFilename(inFilename string) string {
	ext := filepath.Ext(inFilename)
	return strings.TrimSuffix(inFilename, ext) + ".vm"
}

func compileFile(path string) {
	jackFile, err := os.Open(path)
	if err != nil {
		Die("cannot open %s: %v", path, err)
	}
	defer jackFile.Close()

	vmFile, err := os.Create(outFilename(path))
	if err != nil {
		Die("cannot create vm file: %v", err)
	}
	defer vmFile.Close()

	engine := NewCompilationEngine(jackFile, vmFile)
	engine.Compile()
}
