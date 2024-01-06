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
		w.write("@SP") // pop y
		w.write("AM=M-1")
		w.write("D=M")
		w.write("A=A-1") // point x
		w.write("M=D+M") // x = y + x

	case "sub":
		w.write("@SP") // pop y
		w.write("AM=M-1")
		w.write("D=M")
		w.write("A=A-1") // point x
		w.write("M=M-D") // x = x - y

	case "neg":
		w.write("@SP")
		w.write("A=M-1") // point y
		w.write("M=-M")  // y = -y

	case "eq":
		w.write("@SP") // pop y
		w.write("AM=M-1")
		w.write("D=M")
		w.write("A=A-1") // point x
		w.write("D=M-D") // D = x - y

		seq := w.seqGen.gen("SET_TRUE")
		endSetTrueLabel := fmt.Sprintf("END_SET_TRUE$%d", seq)

		// set false
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=0")   // x = false
		w.write(fmt.Sprintf("@%s", endSetTrueLabel))
		w.write("D; JNE")

		// set true
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=-1")  // x = true
		w.write(fmt.Sprintf("(%s)", endSetTrueLabel))

	case "gt":
		w.write("@SP") // pop y
		w.write("AM=M-1")
		w.write("D=M")
		w.write("A=A-1") // point x
		w.write("D=M-D") // D = x - y

		seq := w.seqGen.gen("SET_TRUE")
		endSetTrueLabel := fmt.Sprintf("END_SET_TRUE$%d", seq)

		// set false
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=0")   // x = false
		w.write(fmt.Sprintf("@%s", endSetTrueLabel))
		w.write("D; JLE")

		// set true
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=-1")  // x = true
		w.write(fmt.Sprintf("(%s)", endSetTrueLabel))

	case "lt":
		w.write("@SP") // pop y
		w.write("AM=M-1")
		w.write("D=M")
		w.write("A=A-1") // point x
		w.write("D=M-D") // D = x - y

		seq := w.seqGen.gen("SET_TRUE")
		endSetTrueLabel := fmt.Sprintf("END_SET_TRUE$%d", seq)

		// set false
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=0")   // x = false
		w.write(fmt.Sprintf("@%s", endSetTrueLabel))
		w.write("D; JGE")

		// set true
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=-1")  // x = true
		w.write(fmt.Sprintf("(%s)", endSetTrueLabel))

	case "and":
		w.write("@SP") // pop y
		w.write("AM=M-1")
		w.write("D=M")
		w.write("A=A-1") // point x
		w.write("M=D&M") // x = y & x

	case "or":
		w.write("@SP") // pop y
		w.write("AM=M-1")
		w.write("D=M")
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
