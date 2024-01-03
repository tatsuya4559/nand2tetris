package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
)

type CommandType int

const (
	A_COMMAND CommandType = iota
	C_COMMAND
	L_COMMAND
)

type Command interface {
	Type() CommandType
	String() string
}

type ACommand struct {
	Symbol string
}

func (c *ACommand) Type() CommandType {
	return A_COMMAND
}

func (c *ACommand) String() string {
	return fmt.Sprintf("@%s", c.Symbol)
}

type CCommand struct {
	Dest string
	Comp string
	Jump string
}

func (c *CCommand) Type() CommandType {
	return C_COMMAND
}

func (c *CCommand) String() string {
	return fmt.Sprintf("%s=%s;%s", c.Dest, c.Comp, c.Jump)
}

type LCommand struct {
	Symbol string
}

func (c *LCommand) Type() CommandType {
	return L_COMMAND
}

func (c *LCommand) String() string {
	return fmt.Sprintf("(%s)", c.Symbol)
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

func (p *Parser) CurrentCommand() Command {
	return p.currentCommand
}

var commentPrefix = []byte("//")

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
	cmd, err := parse(word)
	if err != nil {
		log.Fatalf("Cannot parse command: %v", err)
	}
	p.currentCommand = cmd
}

func parse(word string) (Command, error) {
	switch {
	case strings.HasPrefix(word, "@"):
		return parseACommand(word)
	case strings.HasPrefix(word, "("):
		return parseLCommand(word)
	default:
		return parseCCommand(word)
	}
}

func parseACommand(word string) (*ACommand, error) {
	// @symbol
	Assert(word[0] == '@', "A-Command must start with '@'")
	cmd := ACommand{
		Symbol: word[1:],
	}
	return &cmd, nil
}

func parseCCommand(word string) (*CCommand, error) {
	// dest=comp; jump
	cmd := CCommand{}
	if i := strings.Index(word, "="); i >= 0 {
		cmd.Dest = strings.TrimSpace(word[:i])
		word = word[i+1:]
	}
	if i := strings.Index(word, ";"); i >= 0 {
		cmd.Jump = strings.TrimSpace(word[i+1:])
		word = word[:i]
	}
	cmd.Comp = strings.TrimSpace(word)
	return &cmd, nil
}

func parseLCommand(word string) (*LCommand, error) {
	// (SYMBOL)
	if word[len(word)-1] != ')' {
		return nil, fmt.Errorf("L-Command format is invalid: %q", word)
	}
	cmd := LCommand{
		Symbol: strings.TrimSpace(word[1 : len(word)-1]),
	}
	return &cmd, nil
}
