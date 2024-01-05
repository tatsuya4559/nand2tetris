package main

import "fmt"

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
	Symbol        string
	SymbolIsDigit bool
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
