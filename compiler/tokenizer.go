package main

import (
	"io"
	"unicode"
	"unicode/utf8"
)

const (
	eof = -1
)

type Tokenizer struct {
	src      []rune
	ch       rune // current character
	offset   int  // character offset
	rdOffset int  // reading offset(position after current character)
	lineNum  int  // current line number
}

func NewTokenizer(input io.Reader) *Tokenizer {
	b, err := io.ReadAll(input)
	if err != nil {
		Die("cannot read src file: %v", err)
	}
	t := Tokenizer{src: bytesToRunes(b), lineNum: 1}
	t.readRune()
	return &t
}

func bytesToRunes(byteSlice []byte) []rune {
	runeSlice := make([]rune, 0, utf8.RuneCount(byteSlice))
	for len(byteSlice) > 0 {
		r, size := utf8.DecodeRune(byteSlice)
		runeSlice = append(runeSlice, r)
		byteSlice = byteSlice[size:]
	}
	return runeSlice
}

// Read the next run into t.ch
func (t *Tokenizer) readRune() {
	if t.rdOffset < len(t.src) {
		if t.ch == '\n' {
			t.lineNum++
		}
		t.ch = t.src[t.rdOffset]
	} else {
		t.ch = eof
	}
	t.offset = t.rdOffset
	t.rdOffset++
}

func (t *Tokenizer) readString() string {
	t.readRune() // consume "
	begin := t.offset
	for t.ch != '"' {
		t.readRune()
	}
	end := t.offset
	str := string(t.src[begin:end])
	t.readRune() // consume "
	return str
}

func (t *Tokenizer) readInt() string {
	begin := t.offset
	for isDigit(t.ch) {
		t.readRune()
	}
	end := t.offset
	return string(t.src[begin:end])
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isLetter(r rune) bool {
	return ('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		r == '_' ||
		isDigit(r)
}

func (t *Tokenizer) readIdentifier() string {
	begin := t.offset
	for isLetter(t.ch) {
		t.readRune()
	}
	end := t.offset
	return string(t.src[begin:end])
}

func (t *Tokenizer) skipWhiteSpaces() {
	for unicode.IsSpace(t.ch) {
		t.readRune()
	}
}

func (t *Tokenizer) NextToken() Token {
	var tok Token

	t.skipWhiteSpaces()

	switch t.ch {
	case eof:
		tok.Kind = TokenEOF
	case '{', '}', '(', ')', '[', ']', '.', ',', ';',
		'+', '-', '*', '/', '&', '|', '<', '>', '=', '~':
		tok.Kind = TokenSymbol
		tok.Literal = string(t.ch)
		t.readRune() // consume symbol
	case '"':
		tok.Kind = TokenString
		tok.Literal = t.readString()
	default:
		if isDigit(t.ch) {
			tok.Kind = TokenInt
			tok.Literal = t.readInt()
		}
		ident := t.readIdentifier()
		if IsKeyword(ident) {
			tok.Kind = TokenKeyword
			tok.Literal = ident
		} else {
			tok.Kind = TokenIdentifier
			tok.Literal = ident
		}
	}

	return tok
}