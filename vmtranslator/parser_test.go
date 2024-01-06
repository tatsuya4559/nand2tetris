package main

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	vmCode := `	
// TEST
// Pushes and adds two constants.

push constant 7
push constant 8
add
`

	parser := NewParser(strings.NewReader(vmCode))

	if !parser.HasMoreCommands() {
		t.Fatal("HasMoreCommands(1) returned false")
	}
	parser.Advance()
	if got := parser.CommandType(); got != C_PUSH {
		t.Fatalf("expect CommandType() to be C_PUSH(1), but got %v", got)
	}
	if got := parser.Arg1(); got != "constant" {
		t.Fatalf("expect Arg1() to be `constant`, but got %v", got)
	}

	if !parser.HasMoreCommands() {
		t.Fatal("HasMoreCommands(2) returned false")
	}
	parser.Advance()
	if got := parser.CommandType(); got != C_PUSH {
		t.Fatalf("expect CommandType() to be C_PUSH(2), but got %v", got)
	}

	if !parser.HasMoreCommands() {
		t.Fatal("HasMoreCommands(3) returned false")
	}
	parser.Advance()
	if got := parser.CommandType(); got != C_ARITHMETIC {
		t.Fatalf("expect CommandType() to be C_ARITHMETIC, but got %v", got)
	}

	if parser.HasMoreCommands() {
		t.Fatal("HasMoreCommands(4) returned true, but there should be no commands left")
	}
}
