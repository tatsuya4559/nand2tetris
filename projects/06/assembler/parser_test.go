package main

import (
	"bufio"
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
			t.Fatalf("Failed to scan at line %d", i+1)
		}
		if got := scanner.Text(); got != want {
			t.Errorf("Want %q, but got %q in line %d", want, got, i+1)
		}
	}
}
