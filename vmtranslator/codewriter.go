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
	out      io.Writer
	seqGen   sequenceGenerator
	filename string
}

func NewCodeWriter(out io.Writer) *CodeWriter {
	return &CodeWriter{
		out:    out,
		seqGen: make(sequenceGenerator),
	}
}

func (w *CodeWriter) SetFilename(filename string) {
	w.filename = filename
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
}

func (w *CodeWriter) writePush(segment string, index int) {
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
		w.write("@R15")
		w.write("A=A+1")
		w.write("D=A")
		w.write(fmt.Sprintf("@%d", index))
		w.write("A=D+A")
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

	// M[SP] = D
	w.write("@SP")
	w.write("A=M")
	w.write("M=D")
	// SP++
	w.write("@SP")
	w.write("M=M+1")
}

func (w *CodeWriter) writePop(segment string, index int) {
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
		w.write("@R15")
		w.write("A=A+1")
		w.write("D=A")
		w.write(fmt.Sprintf("@%d", index))
		w.write("D=D+A")
	case "constant":
		// We cannot save poped value into constant segment.
		// So we discard it when typ is C_POP.
		return
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
		return
	case "temp":
		// temp segment is R5 ~ R12
		if index < 0 || 7 < index {
			Die("Segmentation Fault: access over temp segment: index=%d", index)
		}
		w.write(fmt.Sprintf("@R%d", index+5))
		w.write("M=D")
		return
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

func (w *CodeWriter) WriteLabel(label string) {
	// TODO: (function$label) の形式にする
	// 現在の関数名がわからないといけない
	w.write(fmt.Sprintf("(%s)", label))
}

func (w *CodeWriter) WriteGoto(label string) {
	// TODO: (function$label) の形式にする
	w.write(fmt.Sprintf("@%s", label))
	w.write("0;JMP")
}

func (w *CodeWriter) WriteIf(label string) {
	// pop
	w.write("@SP")
	w.write("AM=M-1")
	w.write("D=M")

	// TODO: (function$label) の形式にする
	w.write(fmt.Sprintf("@%s", label))
	w.write("D;JNE")
}
