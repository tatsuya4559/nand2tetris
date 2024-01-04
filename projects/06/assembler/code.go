package main

import (
	"fmt"
	"strconv"
)

const (
	_JGT uint16 = 1 << iota
	_JEQ
	_JLT
	_M
	_D
	_A

	_C_INSTRUCTION_MARKER = 0b111 << 13
)

var (
	jumpMnemonicToBinary = map[string]uint16{
		"":    0,
		"JGT": _JGT,
		"JEQ": _JEQ,
		"JGE": _JGT | _JEQ,
		"JLT": _JLT,
		"JNE": _JLT | _JGT,
		"JLE": _JLT | _JEQ,
		"JMP": _JLT | _JEQ | _JGT,
	}
	destMnemonicToBinary = map[string]uint16{
		"":    0,
		"A":   _A,
		"D":   _D,
		"M":   _M,
		"AD":  _A | _D,
		"AM":  _A | _M,
		"MD":  _M | _D,
		"AMD": _A | _M | _D,
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

func CommandToBinaryCode(c Command) (uint16, error) {
	switch cmd := c.(type) {
	case *ACommand:
		addr, err := strconv.ParseUint(cmd.Symbol, 10, 16)
		if err != nil {
			return 0, fmt.Errorf("cannot convert symbol(%s) into uint16: %w", cmd.Symbol, err)
		}
		bin := uint16(addr)
		return bin, nil
	case *CCommand:
		comp := compMnemonicToBinary[cmd.Comp]
		dest := destMnemonicToBinary[cmd.Dest]
		jump := jumpMnemonicToBinary[cmd.Jump]
		return comp | dest | jump | _C_INSTRUCTION_MARKER, nil
	default:
		panic(fmt.Sprintf("Command %v does not have binary representation", c))
	}
}
