package main

import (
	"fmt"
	"io"
)

type Segment string

const (
	SegConst   Segment = "constant"
	SegLocal   Segment = "local"
	SegArg     Segment = "argument"
	SegStatic  Segment = "static"
	SegThis    Segment = "this"
	SegThat    Segment = "that"
	SegPointer Segment = "pointer"
	SegTemp    Segment = "temp"
)

type ArithmeticCommand string

const (
	CmdAdd ArithmeticCommand = "add"
	CmdSub ArithmeticCommand = "sub"
	CmdNeg ArithmeticCommand = "neg"
	CmdEq  ArithmeticCommand = "eq"
	CmdGt  ArithmeticCommand = "gt"
	CmdLt  ArithmeticCommand = "lt"
	CmdAnd ArithmeticCommand = "and"
	CmdOr  ArithmeticCommand = "or"
	CmdNot ArithmeticCommand = "not"
)

type VMWriter struct {
	out io.Writer
}

func NewVMWriter(out io.Writer) *VMWriter {
	return &VMWriter{out: out}
}

func (w *VMWriter) writef(format string, args ...any) {
	fmt.Fprintf(w.out, format, args...)
	fmt.Fprint(w.out, "\n")
}

func (w *VMWriter) WritePush(seg Segment, index int) {
	w.writef("push %s %d", seg, index)
}

func (w *VMWriter) WritePop(seg Segment, index int) {
	w.writef("pop %s %d", seg, index)
}

func (w *VMWriter) WriteArithmeric(cmd ArithmeticCommand) {
	w.writef("%s", cmd)
}

func (w *VMWriter) WriteLabel(label string) {
	w.writef("label %s", label)
}

func (w *VMWriter) WriteGoto(label string) {
	w.writef("goto %s", label)
}

func (w *VMWriter) WriteIf(label string) {
	w.writef("if-goto %s", label)
}

func (w *VMWriter) WriteCall(name string, nArgs int) {
	w.writef("call %s %d", name, nArgs)
}

func (w *VMWriter) WriteFunction(name string, nLocals int) {
	w.writef("function %s %d", name, nLocals)
}

func (w *VMWriter) WriteReturn() {
	w.writef("return")
}
