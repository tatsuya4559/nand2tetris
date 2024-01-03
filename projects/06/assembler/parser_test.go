package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestScanCommand(t *testing.T) {
	input := `// Computes R2 = max(R0, R1)  (R0,R1,R2 refer to RAM[0],RAM[1],RAM[2])

   // D = R0 - R1
   @R0
   D=M
   @R1
   D=D-M
   // If (D > 0) goto ITSR0
   @ITSR0
   D;JGT
   // Its R1
   @R1
   D=M
   @R2
   M=D
   @END // goto end
   0;JMP

(ITSR0)
   @R0             
   D=M
   @R2
   M=D
(END)
   @END
   0;JMP`

	wants := []string{
		"@R0",
		"D=M",
		"@R1",
		"D=D-M",
		"@ITSR0",
		"D;JGT",
		"@R1",
		"D=M",
		"@R2",
		"M=D",
		"@END",
		"0;JMP",
		"(ITSR0)",
		"@R0",
		"D=M",
		"@R2",
		"M=D",
		"(END)",
		"@END",
		"0;JMP",
	}

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(scanCommand)

	for i, want := range wants {
		if !scanner.Scan() {
			t.Fatalf("failed to scan at line %d", i+1)
		}
		if got := scanner.Text(); got != want {
			t.Errorf("want %q, but got %q in line %d", want, got, i+1)
		}
	}
}

func TestParser(t *testing.T) {
	input := `// Computes R2 = max(R0, R1)  (R0,R1,R2 refer to RAM[0],RAM[1],RAM[2])

   // D = R0 - R1
   @R0
   D=M
   @R1
   D=D-M
   // If (D > 0) goto ITSR0
   @ITSR0
   D;JGT
   // Its R1
   @R1
   D=M
   @R2
   M=D
   @END // goto end
   0;JMP

( ITSR0 )
   @R0             
   D  = M
   @R2
   M=D
(END)
   @END
   0; JMP`

	wants := []string{
		"@R0",
		"D=M;",
		"@R1",
		"D=D-M;",
		"@ITSR0",
		"=D;JGT",
		"@R1",
		"D=M;",
		"@R2",
		"M=D;",
		"@END",
		"=0;JMP",
		"(ITSR0)",
		"@R0",
		"D=M;",
		"@R2",
		"M=D;",
		"(END)",
		"@END",
		"=0;JMP",
	}

	parser := NewParser(strings.NewReader(input))

	for i, want := range wants {
		parser.Advance()
		if got := parser.CurrentCommand().String(); got != want {
			t.Errorf("want %q, but got %q in line %d", want, got, i+1)
		}
	}
}

func TestIsValidSymbol(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"foo", true},
		{"foo1._$:", true},
		{"1foo", false},
		{"123", true},
		{"42.195", false},
		{"-1", false},
		{"0", true},
		{"01", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("isValidSymbol(%q)", tt.input), func(t *testing.T) {
			got := isValidSymbol(tt.input)
			if got != tt.want {
				t.Errorf("want %v, but got %v", tt.want, got)
			}
		})
	}
}
