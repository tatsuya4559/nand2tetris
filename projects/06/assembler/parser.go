package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strings"
)

var (
	commentPrefix = []byte("//")
)

type CommandType int

const (
	A_COMMAND CommandType = iota
	C_COMMAND
	L_COMMAND
)

type Command interface {
	Type() CommandType
}

type ACommand struct {
	Symbol string
}

func (c *ACommand) Type() CommandType {
	return A_COMMAND
}

type CCommand struct {
	Dest string
	Comp string
	Jump string
}

func (c *CCommand) Type() CommandType {
	return C_COMMAND
}

type LCommand struct {
	Symbol string
}

func (c *LCommand) Type() CommandType {
	return L_COMMAND
}

type Parser struct {
	scanner        *bufio.Scanner
	buffer         *strings.Builder
	currentCommand Command
	eof            bool
}

func NewParser(r io.Reader) *Parser {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanCommand)
	buffer := &strings.Builder{}
	return &Parser{
		scanner: scanner,
		buffer:  buffer,
	}
}

func (p *Parser) HasMoreCommand() bool {
	return !p.eof
}

// scanCommand scans each asm command ignoring whitespaces, newlines and comments.
// It returns empty string for blank lines and comment only lines.
func scanCommand(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for {
		advance, token, err = bufio.ScanLines(data[start:], atEOF)
		if token == nil || err != nil {
			return
		}
		// If token contains comment, remove it.
		if i := bytes.Index(token, commentPrefix); i >= 0 {
			token = token[:i]
		}
		token = bytes.TrimSpace(token)
		if len(token) > 0 {
			return start + advance, token, nil
		}
		start = start + advance
	}
}

func (p *Parser) Advance() {
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			log.Fatalf("Failed to scan asm file: %v", err)
		}
		p.eof = true
		return
	}

	word := p.scanner.Text()
	log.Printf("%q", word)

	// cmd, err := parse(word)
	// if err != nil {
	// 	log.Fatalf("Cannot parse command: %q", word)
	// }
	// p.currentCommand = cmd
}

// func parse(word string) (Command, error) {
// 	switch {
// 	case strings.HasPrefix(word, "@"):
// 		return parseACommand(word)
// 	case strings.HasPrefix(word, "("):
// 		return parseLCommand(word)
// 	default:
// 		return parseCCommand(word)
// 	}
// }

// func parseACommand(word) (Command, error) {
// 	word
// }
