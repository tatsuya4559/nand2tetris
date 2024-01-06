package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type CommandType int

const (
	C_ARITHMETIC CommandType = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

type Parser struct {
	scanner        *bufio.Scanner
	currentCommand []string
	nextCommand    []string
	isEOF          bool
}

func NewParser(in io.Reader) *Parser {
	scanner := bufio.NewScanner(in)
	parser := Parser{scanner: scanner}
	parser.Advance() // load nextCommand
	return &parser
}

func (p *Parser) HasMoreCommands() bool {
	return !p.isEOF
}

func (p *Parser) Advance() {
	p.currentCommand = p.nextCommand
	for p.scanner.Scan() {
		line := p.scanner.Text()
		line = trimComment(line)

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		p.nextCommand = fields
		return
	}
	p.nextCommand = nil
	p.isEOF = true
}

func trimComment(line string) string {
	if i := strings.Index(line, "//"); i >= 0 {
		return line[:i]
	}
	return line
}

func (p *Parser) CommandType() CommandType {
	switch p.currentCommand[0] {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return C_ARITHMETIC
	case "push":
		return C_PUSH
	case "pop":
		return C_POP
	case "label":
		return C_LABEL
	case "goto":
		return C_GOTO
	case "if-goto":
		return C_IF
	case "function":
		return C_FUNCTION
	case "return":
		return C_RETURN
	case "call":
		return C_CALL
	default:
		Die("unknown command: %s", p.currentCommand[0])
		return -1 // unreachable
	}
}

func (p *Parser) Arg1() string {
	commandType := p.CommandType()
	Assert(commandType != C_RETURN)

	if commandType == C_ARITHMETIC {
		return p.currentCommand[0]
	}
	return p.currentCommand[1]
}

func (p *Parser) Arg2() int {
	commandType := p.CommandType()
	Assert(commandType == C_PUSH ||
		commandType == C_POP ||
		commandType == C_FUNCTION ||
		commandType == C_CALL)

	i, err := strconv.Atoi(p.currentCommand[2])
	if err != nil {
		Die("second argument of command %q was not an integer", p.currentCommand[0])
	}
	return i
}
