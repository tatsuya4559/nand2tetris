package main

import (
	"fmt"
	"io"
)

type sequenceGenerator map[string]int

func (s sequenceGenerator) gen(key string) int {
	result := s[key]
	s[key]++
	return result
}

type CodeWriter struct {
	out    io.Writer
	seqGen sequenceGenerator
}

func NewCodeWriter(out io.Writer) *CodeWriter {
	return &CodeWriter{
		out:    out,
		seqGen: make(sequenceGenerator),
	}
}

func (w *CodeWriter) write(s string) {
	io.WriteString(w.out, s)
	io.WriteString(w.out, "\n")
}

// writePop write asm that means pop stack to D register
func (w *CodeWriter) writePop() {
	// SP--
	w.write("@SP")
	w.write("AM=M-1")
	// D = M[SP]
	w.write("D=M")
}

func (w *CodeWriter) WriteArithmetic(command string) {
	// Comments assume following initial state.
	//  stack
	// +-----+
	// | ... |
	// |  x  |
	// |  y  |
	// |     | <- SP
	switch command {
	case "add":
		w.writePop()     // pop y
		w.write("A=A-1") // point x
		w.write("M=D+M") // x = y + x

	case "sub":
		w.writePop()     // pop y
		w.write("A=A-1") // point x
		w.write("M=M-D") // x = x - y

	case "neg":
		w.write("@SP")
		w.write("A=M-1") // point y
		w.write("M=-M")  // y = -y

	case "eq":
		w.writePop()     // pop y
		w.write("A=A-1") // point x
		w.write("D=M-D") // D = x - y

		seq := w.seqGen.gen("SET_TRUE")
		setTrueLabel := fmt.Sprintf("SET_TRUE$%d", seq)
		endSetTrueLabel := fmt.Sprintf("END_SET_TRUE$%d", seq)

		w.write(fmt.Sprintf("@%s", setTrueLabel))
		w.write("D; JEQ")

		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=0")   // x = false
		w.write(fmt.Sprintf("@%s", endSetTrueLabel))
		w.write("0; JMP")

		w.write(fmt.Sprintf("(%s)", setTrueLabel))
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=-1")  // x = true
		w.write(fmt.Sprintf("(%s)", endSetTrueLabel))

	case "gt":
		w.writePop()     // pop y
		w.write("A=A-1") // point x
		w.write("D=M-D") // D = x - y

		seq := w.seqGen.gen("SET_TRUE")
		setTrueLabel := fmt.Sprintf("SET_TRUE$%d", seq)
		endSetTrueLabel := fmt.Sprintf("END_SET_TRUE$%d", seq)

		w.write(fmt.Sprintf("@%s", setTrueLabel))
		w.write("D; JGT")

		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=0")   // x = false
		w.write(fmt.Sprintf("@%s", endSetTrueLabel))
		w.write("0; JMP")

		w.write(fmt.Sprintf("(%s)", setTrueLabel))
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=-1")  // x = true
		w.write(fmt.Sprintf("(%s)", endSetTrueLabel))

	case "lt":
		w.writePop()     // pop y
		w.write("A=A-1") // point x
		w.write("D=M-D") // D = x - y

		seq := w.seqGen.gen("SET_TRUE")
		setTrueLabel := fmt.Sprintf("SET_TRUE$%d", seq)
		endSetTrueLabel := fmt.Sprintf("END_SET_TRUE$%d", seq)

		w.write(fmt.Sprintf("@%s", setTrueLabel))
		w.write("D; JLT")

		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=0")   // x = false
		w.write(fmt.Sprintf("@%s", endSetTrueLabel))
		w.write("0; JMP")

		w.write(fmt.Sprintf("(%s)", setTrueLabel))
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=-1")  // x = true
		w.write(fmt.Sprintf("(%s)", endSetTrueLabel))

	case "and":
		w.writePop()     // pop y
		w.write("A=A-1") // point x
		w.write("M=D&M") // x = y & x

	case "or":
		w.writePop()     // pop y
		w.write("A=A-1") // point x
		w.write("M=D|M") // x = y | x

	case "not":
		w.write("@SP")
		w.write("A=M-1") // point y
		w.write("M=!M")  // y = !y
	}
}

func (w *CodeWriter) WritePushPop(typ CommandType, segment string, index int) {
	if segment == "constant" {
		// D = index
		w.write(fmt.Sprintf("@%d", index))
		w.write("D=A")
	}

	if typ == C_PUSH {
		// M[SP] = D
		w.write("@SP")
		w.write("A=M")
		w.write("M=D")
		// SP++
		w.write("@SP")
		w.write("M=M+1")
	}

	// push foo ->
	// D=whatToPush
	// @SP
	// M=D
}
