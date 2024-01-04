package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

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

var (
	commentPrefix = []byte("//")
	whiteSpaces   = []string{" ", "\t"}
)

func removeWhiteSpaces(s []byte) []byte {
	for _, space := range whiteSpaces {
		s = bytes.ReplaceAll(s, []byte(space), []byte(""))
	}
	return s
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
		token = removeWhiteSpaces(token)
		if len(token) > 0 {
			return start + advance, token, nil
		}
		start = start + advance
	}
}

func (p *Parser) Advance() (scanned bool, err error) {
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			return false, fmt.Errorf("failed to scan asm file: %w", err)
		}
		p.eof = true
		return false, nil
	}

	word := p.scanner.Text()
	cmd, err := parse(word)
	if err != nil {
		return false, err
	}
	p.currentCommand = cmd
	return true, nil
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

// isValidDecimal check symbol is a valid decimal.
// It must not be negative.
func isValidDecimal(symbol string) bool {
	for i, r := range symbol {
		if i == 0 && r == '0' {
			return false
		}
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isSymbol(r rune) bool {
	return unicode.IsDigit(r) ||
		('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		r == '_' ||
		r == '.' ||
		r == '$' ||
		r == ':'
}

func isValidName(symbol string) bool {
	for i, r := range symbol {
		if len(symbol) > 1 && i == 0 && unicode.IsDigit(r) {
			return false
		}
		if !isSymbol(r) {
			return false
		}
	}
	return true
}

func isValidSymbol(symbol string) bool {
	return isValidDecimal(symbol) || isValidName(symbol)
}

var (
	compMnemonics = NewSet("0", "1", "-1", "D", "A", "!D", "!A", "-D", "-A",
		"D+1", "A+1", "D-1", "A-1", "D+A", "D-A", "A-D", "D&A", "D|A",
		"M", "!M", "-M", "M+1", "M-1", "D+M", "D-M", "M-D", "D&M", "D|M")
	destMnemonics = NewSet("", "A", "D", "M", "AD", "AM", "MD", "AMD")
	jumpMnemonics = NewSet("", "JGT", "JEQ", "JGE", "JLT", "JNE", "JLE", "JMP")
)

func parseACommand(word string) (*ACommand, error) {
	// @symbol
	Assert(word[0] == '@', "A-Command must start with '@'")
	symbol := word[1:]
	if !isValidSymbol(symbol) {
		return nil, NewParseError("invalid symbol", symbol)
	}
	cmd := ACommand{Symbol: symbol}
	return &cmd, nil
}

func parseCCommand(word string) (*CCommand, error) {
	// dest=comp; jump
	cmd := CCommand{}
	if i := strings.Index(word, "="); i >= 0 {
		dest := word[:i]
		if !destMnemonics.Contains(dest) {
			return nil, NewParseError("invalid dest mnemonic", dest)
		}
		cmd.Dest = dest
		word = word[i+1:]
	}
	if i := strings.Index(word, ";"); i >= 0 {
		jump := word[i+1:]
		if !jumpMnemonics.Contains(jump) {
			return nil, NewParseError("invalid jump mnemonic", jump)
		}
		cmd.Jump = jump
		word = word[:i]
	}
	if !compMnemonics.Contains(word) {
		return nil, NewParseError("invalid comp mnemonic", word)
	}
	cmd.Comp = word
	return &cmd, nil
}

func parseLCommand(word string) (*LCommand, error) {
	// (SYMBOL)
	if word[len(word)-1] != ')' {
		return nil, NewParseError("closing paren is not found", word)
	}
	symbol := word[1 : len(word)-1]
	if !isValidSymbol(symbol) {
		return nil, NewParseError("invalid symbol", symbol)
	}
	cmd := LCommand{Symbol: symbol}
	return &cmd, nil
}

type ParseError struct {
	message string
	word    string
	err     error
}

func NewParseError(message, word string) *ParseError {
	return &ParseError{
		message: message,
		word:    word,
	}
}

func (e *ParseError) Error() string {
	if e.err == nil {
		return fmt.Sprintf("parse error: %s at %q", e.message, e.word)
	}
	return fmt.Sprintf("parse error: %s at %q: %v", e.message, e.word, e.err)
}

func (e *ParseError) Unwrap() error {
	return e.err
}
