package main

import (
	"fmt"
	"strconv"
)

const (
	JGT uint16 = 1 << iota
	JEQ
	JLT
	M
	D
	A

	C_INSTRUCTION_MARKER = 0b111 << 13
)

var (
	jumpMnemonicToBinary = map[string]uint16{
		"":    0,
		"JGT": JGT,
		"JEQ": JEQ,
		"JGE": JGT | JEQ,
		"JLT": JLT,
		"JNE": JLT | JGT,
		"JLE": JLT | JEQ,
		"JMP": JLT | JEQ | JGT,
	}
	destMnemonicToBinary = map[string]uint16{
		"":    0,
		"A":   A,
		"D":   D,
		"M":   M,
		"AD":  A | D,
		"AM":  A | M,
		"MD":  M | D,
		"AMD": A | M | D,
	}
	compMnemonicToBinary = map[string]uint16{
		"0":   0b0_101_010 << 6,
		"1":   0b0_111_111 << 6,
		"-1":  0b0_111_010 << 6,
		"D":   0b0_001_100 << 6,
		"A":   0b0_110_000 << 6,
		"!D":  0b0_001_101 << 6,
		"!A":  0b0_110_001 << 6,
		"-D":  0b0_001_111 << 6,
		"-A":  0b0_110_011 << 6,
		"D+1": 0b0_011_111 << 6,
		"A+1": 0b0_110_111 << 6,
		"D-1": 0b0_001_110 << 6,
		"A-1": 0b0_110_010 << 6,
		"D+A": 0b0_000_010 << 6,
		"D-A": 0b0_010_011 << 6,
		"A-D": 0b0_000_111 << 6,
		"D&A": 0b0_000_000 << 6,
		"D|A": 0b0_010_101 << 6,
		"M":   0b1_110_000 << 6,
		"!M":  0b1_110_001 << 6,
		"-M":  0b1_110_011 << 6,
		"M+1": 0b1_110_111 << 6,
		"M-1": 0b1_110_010 << 6,
		"D+M": 0b1_000_010 << 6,
		"D-M": 0b1_010_011 << 6,
		"M-D": 0b1_000_111 << 6,
		"D&M": 0b1_000_000 << 6,
		"D|M": 0b1_010_101 << 6,
	}
)

func ConvertACommand(c *ACommand) (uint16, error) {
	addr, err := strconv.ParseUint(c.Symbol, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("cannot convert symbol(%s) into uint16: %w", c.Symbol, err)
	}
	bin := uint16(addr)
	return bin, nil
}

func ConvertCCommand(c *CCommand) (uint16, error) {
	comp := compMnemonicToBinary[c.Comp]
	dest := destMnemonicToBinary[c.Dest]
	jump := jumpMnemonicToBinary[c.Jump]
	return comp | dest | jump | C_INSTRUCTION_MARKER, nil
}
