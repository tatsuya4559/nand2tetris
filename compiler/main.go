package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
			tokenizeFile(jackFile)
		}
	} else {
		tokenizeFile(path)
	}
}

func tokenizeFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		Die("cannot open %s: %v", path, err)
	}
	defer file.Close()

	engine := NewCompilationEngine(file)
	ast := engine.Compile()
	fmt.Printf("%+v\n", ast)
}
