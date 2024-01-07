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

func (w *CodeWriter) write(s string) {
	io.WriteString(w.out, s)
	io.WriteString(w.out, "\n")
}

func (w *CodeWriter) genSequencialLabel(prefix string) string {
	seq := w.seqGen.gen(prefix)
	return fmt.Sprintf("%s_%d", prefix, seq)
}

func (w *CodeWriter) WriteArithmetic(command string) {
	w.write(fmt.Sprintf("// %s", command))

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

		endSetTrueLabel := w.genSequencialLabel("END_SET_TRUE")

		// set false
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=0")   // x = false
		w.write(fmt.Sprintf("@%s", endSetTrueLabel))
		w.write("D;JNE")

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

		endSetTrueLabel := w.genSequencialLabel("END_SET_TRUE")

		// set false
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=0")   // x = false
		w.write(fmt.Sprintf("@%s", endSetTrueLabel))
		w.write("D;JLE")

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

		endSetTrueLabel := w.genSequencialLabel("END_SET_TRUE")

		// set false
		w.write("@SP")
		w.write("A=M-1") // point x
		w.write("M=0")   // x = false
		w.write(fmt.Sprintf("@%s", endSetTrueLabel))
		w.write("D;JGE")

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

	w.write("")
}

func (w *CodeWriter) writePush(segment string, index int) {
	w.write(fmt.Sprintf("// push %s %d", segment, index))

	switch segment {
	case "argument":
		w.write("@ARG")
		w.write("D=M")
		w.write(fmt.Sprintf("@%d", index))
		w.write("A=D+A")
		w.write("D=M")
	case "local":
		w.write("@LCL")
		w.write("D=M")
		w.write(fmt.Sprintf("@%d", index))
		w.write("A=D+A")
		w.write("D=M")
	case "static":
		w.write(fmt.Sprintf("@%s.static_%d", w.currentFile, index))
		w.write("D=M")
	case "constant":
		w.write(fmt.Sprintf("@%d", index))
		w.write("D=A")
	case "this":
		w.write("@THIS")
		w.write("D=M")
		w.write(fmt.Sprintf("@%d", index))
		w.write("A=D+A")
		w.write("D=M")
	case "that":
		w.write("@THAT")
		w.write("D=M")
		w.write(fmt.Sprintf("@%d", index))
		w.write("A=D+A")
		w.write("D=M")
	case "pointer":
		if index == 0 {
			w.write("@THIS")
		} else if index == 1 {
			w.write("@THAT")
		} else {
			Die("Segmentation Fault: access over pointer segment: index=%d", index)
		}
		w.write("D=M")
	case "temp":
		// temp segment is R5 ~ R12
		if index < 0 || 7 < index {
			Die("Segmentation Fault: access over temp segment: index=%d", index)
		}
		w.write(fmt.Sprintf("@R%d", index+5))
		w.write("D=M")
	default:
		Die("Unknown segment: %s", segment)
	}

	w.writePushD()

	w.write("")
}

// writePushD writes asm which means push D-Register.
func (w *CodeWriter) writePushD() {
	// M[SP] = D
	w.write("@SP")
	w.write("A=M")
	w.write("M=D")
	// SP++
	w.write("@SP")
	w.write("M=M+1")
}

func (w *CodeWriter) writePop(segment string, index int) {
	w.write(fmt.Sprintf("// pop %s %d", segment, index))

	// SP--
	w.write("@SP")
	w.write("AM=M-1")
	// D = M[SP]
	w.write("D=M")
	// R13 = D
	w.write("@R13")
	w.write("M=D")

	// D and R13 hold poped value at this point.

	switch segment {
	case "argument":
		w.write("@ARG")
		w.write("D=M")
		w.write(fmt.Sprintf("@%d", index))
		w.write("D=D+A")
	case "local":
		w.write("@LCL")
		w.write("D=M")
		w.write(fmt.Sprintf("@%d", index))
		w.write("D=D+A")
	case "static":
		w.write(fmt.Sprintf("@%s.static_%d", w.currentFile, index))
		w.write("D=A")
	case "constant":
		// We cannot save poped value into constant segment.
		// So we discard it when typ is C_POP.
		goto END
	case "this":
		w.write("@THIS")
		w.write("D=M")
		w.write(fmt.Sprintf("@%d", index))
		w.write("D=D+A")
	case "that":
		w.write("@THAT")
		w.write("D=M")
		w.write(fmt.Sprintf("@%d", index))
		w.write("D=D+A")
	case "pointer":
		if index == 0 {
			w.write("@THIS")
		} else if index == 1 {
			w.write("@THAT")
		} else {
			Die("Segmentation Fault: access over pointer segment: index=%d", index)
		}
		w.write("M=D")
		goto END
	case "temp":
		// temp segment is R5 ~ R12
		if index < 0 || 7 < index {
			Die("Segmentation Fault: access over temp segment: index=%d", index)
		}
		w.write(fmt.Sprintf("@R%d", index+5))
		w.write("M=D")
		goto END
	default:
		Die("Unknown segment: %s", segment)
	}

	// D holds destination address at this point.
	// First we save the address to R14.
	w.write("@R14")
	w.write("M=D")
	// Then copy the value in R13 to where R14 points.
	w.write("@R13")
	w.write("D=M")
	w.write("@R14")
	w.write("A=M")
	w.write("M=D")

END:
	w.write("")
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
	w.write(fmt.Sprintf("// label %s", label))
	w.write(fmt.Sprintf("(%s)", w.qualifyLabel(label)))
	w.write("")
}

func (w *CodeWriter) WriteGoto(label string) {
	w.write(fmt.Sprintf("// goto %s", label))
	w.write(fmt.Sprintf("@%s", w.qualifyLabel(label)))
	w.write("0;JMP")
	w.write("")
}

func (w *CodeWriter) WriteIf(label string) {
	w.write(fmt.Sprintf("// if-goto %s", label))
	// pop
	w.write("@SP")
	w.write("AM=M-1")
	w.write("D=M")

	w.write(fmt.Sprintf("@%s", w.qualifyLabel(label)))
	w.write("D;JNE")

	w.write("")
}

func (w *CodeWriter) WriteCall(funcName string, nArgs int) {
	w.write(fmt.Sprintf("// call %s %d", funcName, nArgs))

	// Push return address
	returnAddressLabel := w.genSequencialLabel("RETURN_ADDR")
	w.write(fmt.Sprintf("@%s", returnAddressLabel))
	w.write("D=A")
	w.writePushD()

	// Push LCL, ARG, THIS, THAT
	for _, label := range []string{"LCL", "ARG", "THIS", "THAT"} {
		w.write(fmt.Sprintf("@%s", label))
		w.write("D=M")
		w.writePushD()
	}

	// Set ARG to SP - nArgs - 5(return-address, LCL, ARG, THIS, THAT)
	w.write(fmt.Sprintf("@%d", nArgs+5))
	w.write("D=A")
	w.write("@SP")
	w.write("D=M-D")
	w.write("@ARG")
	w.write("M=D")

	// Set LCL to SP
	w.write("@SP")
	w.write("D=M")
	w.write("@LCL")
	w.write("M=D")

	// Goto function
	w.write(fmt.Sprintf("@%s", funcName))
	w.write("0;JMP")

	// Set return address label
	w.write(fmt.Sprintf("(%s) // back from %s to %s", returnAddressLabel, funcName, w.currentFunction))

	w.write("")
}

func (w *CodeWriter) WriteFunction(funcName string, nLocals int) {
	w.currentFunction = funcName

	w.write(fmt.Sprintf("// function %s %d", funcName, nLocals))
	w.write(fmt.Sprintf("(%s) // {", funcName))

	// Initialize local variables
	w.write("D=0")
	for i := 0; i < nLocals; i++ {
		w.writePushD()
	}

	w.write("")
}

func (w *CodeWriter) WriteReturn() {
	w.write(fmt.Sprintf("// return (from %s)", w.currentFunction))

	// Use R15 for saving return address.
	// We need to get return address first because
	// when nargs == 0, return address will be lost
	// by *ARG = pop() operation.
	w.write("@5")
	w.write("D=A")
	w.write("@LCL")
	w.write("A=M-D")
	w.write("D=M")
	w.write("@R15")
	w.write("M=D")

	// Set return value
	w.writePop("argument", 0)
	w.write("@ARG")
	w.write("D=M")
	w.write("@SP")
	w.write("M=D+1")

	// Recover LCL, ARG, THIS, THAT
	for i, label := range []string{"THAT", "THIS", "ARG", "LCL"} {
		w.write(fmt.Sprintf("@%d", i+1))
		w.write("D=A")
		w.write("@LCL")
		w.write("A=M-D")
		w.write("D=M")
		w.write(fmt.Sprintf("@%s", label))
		w.write("M=D")
	}

	// Goto return address.
	w.write("@R15")
	w.write("A=M")
	w.write("0;JMP")

	w.write("// }")
	w.write("")
}

func (w *CodeWriter) WriteInit() {
	// Set SP to RAM[256]
	w.write("@256")
	w.write("D=A")
	w.write("@SP")
	w.write("M=D")

	// Jump to Sys.init
	w.currentFunction = "Sys.init"
	w.WriteCall("Sys.init", 0)
}
