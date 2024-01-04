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
	currentCommand Command
	eof            bool
}

func NewParser(r io.Reader) *Parser {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanCommand)
	return &Parser{scanner: scanner}
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

// Parse scans a command and parse it.
// If successfully parsed it returns true, otherwise returns false.
func (p *Parser) Parse() bool {
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			Die("failed to scan asm file: %v", err)
		}
		p.eof = true
		p.currentCommand = nil
		return false
	}

	word := p.scanner.Text()
	cmd, err := parse(word)
	if err != nil {
		Die(err.Error())
	}

	p.currentCommand = cmd
	return true
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

func isLetter(r rune) bool {
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
		if !isLetter(r) {
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

// parseACommand parses `@symbol`
func parseACommand(word string) (*ACommand, error) {
	symbol := word[1:] // remove @
	if !isValidSymbol(symbol) {
		return nil, NewParseError("invalid symbol", symbol)
	}
	cmd := ACommand{Symbol: symbol}
	return &cmd, nil
}

// parseCCommand parses `dest=comp; jump`
func parseCCommand(word string) (*CCommand, error) {
	cmd := CCommand{}
	if i := strings.Index(word, "="); i >= 0 {
		dest := word[:i]
		if !destMnemonics.Contains(dest) {
			return nil, NewParseError("unknown dest mnemonic", dest)
		}
		cmd.Dest = dest
		word = word[i+1:]
	}
	if i := strings.Index(word, ";"); i >= 0 {
		jump := word[i+1:]
		if !jumpMnemonics.Contains(jump) {
			return nil, NewParseError("unknown jump mnemonic", jump)
		}
		cmd.Jump = jump
		word = word[:i]
	}
	if !compMnemonics.Contains(word) {
		return nil, NewParseError("unknown comp mnemonic", word)
	}
	cmd.Comp = word
	return &cmd, nil
}

// parseLCommand parses `(SYMBOL)`
func parseLCommand(word string) (*LCommand, error) {
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
}

func NewParseError(message, word string) *ParseError {
	return &ParseError{
		message: message,
		word:    word,
	}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error: %s at %q", e.message, e.word)
}
