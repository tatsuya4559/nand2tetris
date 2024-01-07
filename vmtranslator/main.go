package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)

	info, err := os.Stat(path)
	if err != nil {
		Die("cannot stat %s", path)
	}

	var asmFilename string

	if info.IsDir() {
		abspath, err := filepath.Abs(path)
		if err != nil {
			Die("cannot get absolute path for %s", info.Name())
		}
		asmFilename = filepath.Join(path, filepath.Base(abspath)+".asm")

		entries, err := os.ReadDir(path)
		if err != nil {
			Die("cannot read dir %s", path)
		}

		out, err := os.Create(asmFilename)
		if err != nil {
			Die("cannot create %s: %v", asmFilename, err)
		}
		codeWriter := NewCodeWriter(out)
		defer out.Close()

		// codeWriter.WriteInit()
		for _, e := range entries {
			name := e.Name()
			if filepath.Ext(name) != ".vm" {
				continue
			}
			translateVM(filepath.Join(abspath, name), codeWriter)
		}
	} else {
		ext := filepath.Ext(path)
		if ext != ".vm" {
			Die("please give .vm file or directory")
		}
		asmFilename = strings.TrimSuffix(path, ext) + ".asm"

		out, err := os.Create(asmFilename)
		if err != nil {
			Die("cannot create %s: %v", asmFilename, err)
		}
		codeWriter := NewCodeWriter(out)
		defer out.Close()

		// codeWriter.WriteInit()
		translateVM(path, codeWriter)
	}
}

func translateVM(vmPath string, w *CodeWriter) {
	in, err := os.Open(vmPath)
	if err != nil {
		Die("cannot open %s", vmPath)
	}
	defer in.Close()

	w.SetFilename(filepath.Base(vmPath))

	p := NewParser(in)
	for p.HasMoreCommands() {
		p.Advance()
		switch typ := p.CommandType(); typ {
		case C_ARITHMETIC:
			w.WriteArithmetic(p.Arg1())
		case C_PUSH, C_POP:
			w.WritePushPop(typ, p.Arg1(), p.Arg2())
		case C_LABEL:
			w.WriteLabel(p.Arg1())
		case C_GOTO:
			w.WriteGoto(p.Arg1())
		case C_IF:
			w.WriteIf(p.Arg1())
		case C_CALL:
			w.WriteCall(p.Arg1(), p.Arg2())
		case C_FUNCTION:
			w.WriteFunction(p.Arg1(), p.Arg2())
		case C_RETURN:
			w.WriteReturn()
		}
	}
}
