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
	out             io.Writer
	seqGen          sequenceGenerator
	currentFile     string
	currentFunction string
}

func NewCodeWriter(out io.Writer) *CodeWriter {
	return &CodeWriter{
		out:    out,
		seqGen: make(sequenceGenerator),
	}
}

func (w *CodeWriter) SetFilename(filename string) {
	w.currentFile = filename
}

func (w *CodeWriter) writef(format string, args ...any) {
	fmt.Fprintf(w.out, format, args...)
	io.WriteString(w.out, "\n")
}

func (w *CodeWriter) genSequencialLabel(prefix string) string {
	seq := w.seqGen.gen(prefix)
	return fmt.Sprintf("%s_%d", prefix, seq)
}

func (w *CodeWriter) WriteArithmetic(command string) {
	w.writef("// %s", command)

	// Comments assume following initial state.
	//  stack
	// +-----+
	// | ... |
	// |  x  |
	// |  y  |
	// |     | <- SP
	switch command {
	case "add":
		w.writef("@SP") // pop y
		w.writef("AM=M-1")
		w.writef("D=M")
		w.writef("A=A-1") // point x
		w.writef("M=D+M") // x = y + x

	case "sub":
		w.writef("@SP") // pop y
		w.writef("AM=M-1")
		w.writef("D=M")
		w.writef("A=A-1") // point x
		w.writef("M=M-D") // x = x - y

	case "neg":
		w.writef("@SP")
		w.writef("A=M-1") // point y
		w.writef("M=-M")  // y = -y

	case "eq":
		w.writef("@SP") // pop y
		w.writef("AM=M-1")
		w.writef("D=M")
		w.writef("A=A-1") // point x
		w.writef("D=M-D") // D = x - y

		endSetTrueLabel := w.genSequencialLabel("END_SET_TRUE")

		// set false
		w.writef("@SP")
		w.writef("A=M-1") // point x
		w.writef("M=0")   // x = false
		w.writef("@%s", endSetTrueLabel)
		w.writef("D;JNE")

		// set true
		w.writef("@SP")
		w.writef("A=M-1") // point x
		w.writef("M=-1")  // x = true
		w.writef("(%s)", endSetTrueLabel)

	case "gt":
		w.writef("@SP") // pop y
		w.writef("AM=M-1")
		w.writef("D=M")
		w.writef("A=A-1") // point x
		w.writef("D=M-D") // D = x - y

		endSetTrueLabel := w.genSequencialLabel("END_SET_TRUE")

		// set false
		w.writef("@SP")
		w.writef("A=M-1") // point x
		w.writef("M=0")   // x = false
		w.writef("@%s", endSetTrueLabel)
		w.writef("D;JLE")

		// set true
		w.writef("@SP")
		w.writef("A=M-1") // point x
		w.writef("M=-1")  // x = true
		w.writef("(%s)", endSetTrueLabel)

	case "lt":
		w.writef("@SP") // pop y
		w.writef("AM=M-1")
		w.writef("D=M")
		w.writef("A=A-1") // point x
		w.writef("D=M-D") // D = x - y

		endSetTrueLabel := w.genSequencialLabel("END_SET_TRUE")

		// set false
		w.writef("@SP")
		w.writef("A=M-1") // point x
		w.writef("M=0")   // x = false
		w.writef("@%s", endSetTrueLabel)
		w.writef("D;JGE")

		// set true
		w.writef("@SP")
		w.writef("A=M-1") // point x
		w.writef("M=-1")  // x = true
		w.writef("(%s)", endSetTrueLabel)

	case "and":
		w.writef("@SP") // pop y
		w.writef("AM=M-1")
		w.writef("D=M")
		w.writef("A=A-1") // point x
		w.writef("M=D&M") // x = y & x

	case "or":
		w.writef("@SP") // pop y
		w.writef("AM=M-1")
		w.writef("D=M")
		w.writef("A=A-1") // point x
		w.writef("M=D|M") // x = y | x

	case "not":
		w.writef("@SP")
		w.writef("A=M-1") // point y
		w.writef("M=!M")  // y = !y
	}

	w.writef("")
}

func (w *CodeWriter) writePush(segment string, index int) {
	w.writef("// push %s %d", segment, index)

	switch segment {
	case "argument":
		w.writef("@ARG")
		w.writef("D=M")
		w.writef("@%d", index)
		w.writef("A=D+A")
		w.writef("D=M")
	case "local":
		w.writef("@LCL")
		w.writef("D=M")
		w.writef("@%d", index)
		w.writef("A=D+A")
		w.writef("D=M")
	case "static":
		w.writef("@%s.static_%d", w.currentFile, index)
		w.writef("D=M")
	case "constant":
		w.writef("@%d", index)
		w.writef("D=A")
	case "this":
		w.writef("@THIS")
		w.writef("D=M")
		w.writef("@%d", index)
		w.writef("A=D+A")
		w.writef("D=M")
	case "that":
		w.writef("@THAT")
		w.writef("D=M")
		w.writef("@%d", index)
		w.writef("A=D+A")
		w.writef("D=M")
	case "pointer":
		if index == 0 {
			w.writef("@THIS")
		} else if index == 1 {
			w.writef("@THAT")
		} else {
			Die("Segmentation Fault: access over pointer segment: index=%d", index)
		}
		w.writef("D=M")
	case "temp":
		// temp segment is R5 ~ R12
		if index < 0 || 7 < index {
			Die("Segmentation Fault: access over temp segment: index=%d", index)
		}
		w.writef("@R%d", index+5)
		w.writef("D=M")
	default:
		Die("Unknown segment: %s", segment)
	}

	w.writePushD()

	w.writef("")
}

// writePushD writes asm which means push D-Register.
func (w *CodeWriter) writePushD() {
	// M[SP] = D
	w.writef("@SP")
	w.writef("A=M")
	w.writef("M=D")
	// SP++
	w.writef("@SP")
	w.writef("M=M+1")
}

func (w *CodeWriter) writePop(segment string, index int) {
	w.writef("// pop %s %d", segment, index)

	// SP--
	w.writef("@SP")
	w.writef("AM=M-1")
	// D = M[SP]
	w.writef("D=M")
	// R13 = D
	w.writef("@R13")
	w.writef("M=D")

	// D and R13 hold poped value at this point.

	switch segment {
	case "argument":
		w.writef("@ARG")
		w.writef("D=M")
		w.writef("@%d", index)
		w.writef("D=D+A")
	case "local":
		w.writef("@LCL")
		w.writef("D=M")
		w.writef("@%d", index)
		w.writef("D=D+A")
	case "static":
		w.writef("@%s.static_%d", w.currentFile, index)
		w.writef("D=A")
	case "constant":
		// We cannot save poped value into constant segment.
		// So we discard it when typ is C_POP.
		goto END
	case "this":
		w.writef("@THIS")
		w.writef("D=M")
		w.writef("@%d", index)
		w.writef("D=D+A")
	case "that":
		w.writef("@THAT")
		w.writef("D=M")
		w.writef("@%d", index)
		w.writef("D=D+A")
	case "pointer":
		if index == 0 {
			w.writef("@THIS")
		} else if index == 1 {
			w.writef("@THAT")
		} else {
			Die("Segmentation Fault: access over pointer segment: index=%d", index)
		}
		w.writef("M=D")
		goto END
	case "temp":
		// temp segment is R5 ~ R12
		if index < 0 || 7 < index {
			Die("Segmentation Fault: access over temp segment: index=%d", index)
		}
		w.writef("@R%d", index+5)
		w.writef("M=D")
		goto END
	default:
		Die("Unknown segment: %s", segment)
	}

	// D holds destination address at this point.
	// First we save the address to R14.
	w.writef("@R14")
	w.writef("M=D")
	// Then copy the value in R13 to where R14 points.
	w.writef("@R13")
	w.writef("D=M")
	w.writef("@R14")
	w.writef("A=M")
	w.writef("M=D")

END:
	w.writef("")
}

func (w *CodeWriter) WritePushPop(typ CommandType, segment string, index int) {
	switch typ {
	case C_PUSH:
		w.writePush(segment, index)
	case C_POP:
		w.writePop(segment, index)
	default:
		Die("Invalid command type for WritePushPop: %v", typ)
	}
}

// qualifyLabel makes label identified in asm file.
func (w *CodeWriter) qualifyLabel(label string) string {
	return fmt.Sprintf("%s.%s$%s", w.currentFile, w.currentFunction, label)
}

func (w *CodeWriter) WriteLabel(label string) {
	w.writef("// label %s", label)
	w.writef("(%s)", w.qualifyLabel(label))
	w.writef("")
}

func (w *CodeWriter) WriteGoto(label string) {
	w.writef("// goto %s", label)
	w.writef("@%s", w.qualifyLabel(label))
	w.writef("0;JMP")
	w.writef("")
}

func (w *CodeWriter) WriteIf(label string) {
	w.writef("// if-goto %s", label)
	// pop
	w.writef("@SP")
	w.writef("AM=M-1")
	w.writef("D=M")

	w.writef("@%s", w.qualifyLabel(label))
	w.writef("D;JNE")

	w.writef("")
}

func (w *CodeWriter) WriteCall(funcName string, nArgs int) {
	w.writef("// call %s %d", funcName, nArgs)

	// Push return address
	returnAddressLabel := w.genSequencialLabel("RETURN_ADDR")
	w.writef("@%s", returnAddressLabel)
	w.writef("D=A")
	w.writePushD()

	// Push LCL, ARG, THIS, THAT
	for _, label := range []string{"LCL", "ARG", "THIS", "THAT"} {
		w.writef("@%s", label)
		w.writef("D=M")
		w.writePushD()
	}

	// Set ARG to SP - nArgs - 5(return-address, LCL, ARG, THIS, THAT)
	w.writef("@%d", nArgs+5)
	w.writef("D=A")
	w.writef("@SP")
	w.writef("D=M-D")
	w.writef("@ARG")
	w.writef("M=D")

	// Set LCL to SP
	w.writef("@SP")
	w.writef("D=M")
	w.writef("@LCL")
	w.writef("M=D")

	// Goto function
	w.writef("@%s", funcName)
	w.writef("0;JMP")

	// Set return address label
	w.writef("(%s) // back from %s to %s", returnAddressLabel, funcName, w.currentFunction)

	w.writef("")
}

func (w *CodeWriter) WriteFunction(funcName string, nLocals int) {
	w.currentFunction = funcName

	w.writef("// function %s %d", funcName, nLocals)
	w.writef("(%s) // {", funcName)

	// Initialize local variables
	w.writef("D=0")
	for i := 0; i < nLocals; i++ {
		w.writePushD()
	}

	w.writef("")
}

func (w *CodeWriter) WriteReturn() {
	w.writef("// return (from %s)", w.currentFunction)

	// Use R15 for saving return address.
	// We need to get return address first because
	// when nargs == 0, return address will be lost
	// by *ARG = pop() operation.
	w.writef("@5")
	w.writef("D=A")
	w.writef("@LCL")
	w.writef("A=M-D")
	w.writef("D=M")
	w.writef("@R15")
	w.writef("M=D")

	// Set return value
	w.writePop("argument", 0)
	w.writef("@ARG")
	w.writef("D=M")
	w.writef("@SP")
	w.writef("M=D+1")

	// Recover LCL, ARG, THIS, THAT
	for i, label := range []string{"THAT", "THIS", "ARG", "LCL"} {
		w.writef("@%d", i+1)
		w.writef("D=A")
		w.writef("@LCL")
		w.writef("A=M-D")
		w.writef("D=M")
		w.writef("@%s", label)
		w.writef("M=D")
	}

	// Goto return address.
	w.writef("@R15")
	w.writef("A=M")
	w.writef("0;JMP")

	w.writef("// }")
	w.writef("")
}

func (w *CodeWriter) WriteInit() {
	// Set SP to RAM[256]
	w.writef("@256")
	w.writef("D=A")
	w.writef("@SP")
	w.writef("M=D")

	// Jump to Sys.init
	w.currentFunction = "Sys.init"
	w.WriteCall("Sys.init", 0)
}
