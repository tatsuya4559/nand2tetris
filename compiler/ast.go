package main

type Expression interface {
	expression() // marker method
}

type BoolLiteralExpression struct {
	Value bool
}

func (e *BoolLiteralExpression) expression() {}

type IntLiteralExpression struct {
	Value int
}

func (e *IntLiteralExpression) expression() {}

type StringLiteralExpression struct {
	Value string
}

func (e *StringLiteralExpression) expression() {}

type ThisLiteralExpression struct {
}

func (e *ThisLiteralExpression) expression() {}

type NullLiteralExpression struct {
}

func (e *NullLiteralExpression) expression() {}

type IdentExpression struct {
	Value string
}

func (e *IdentExpression) expression() {}

type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (e *PrefixExpression) expression() {}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (e *InfixExpression) expression() {}

type DotAccessExpression struct {
	Left  Expression
	Right Expression
}

func (e *DotAccessExpression) expression() {}

type SubroutineCallExpression struct {
	Func Expression
	Args []Expression
}

func (e *SubroutineCallExpression) expression() {}

var (
	NULL  = &NullLiteralExpression{}
	THIS  = &ThisLiteralExpression{}
	TRUE  = &BoolLiteralExpression{Value: true}
	FALSE = &BoolLiteralExpression{Value: false}
)
