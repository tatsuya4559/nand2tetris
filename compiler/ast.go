package main

type ClassDeclaration struct {
	Name        string
	Vars        []*ClassVarDeclaration
	Subroutines []*SubroutineDeclaration
}

type ClassVarDeclaration struct {
	StorageClass string // static or field
	Type         string
	Name         string
}

type SubroutineDeclaration struct {
	Kind       string // constructor or function or method
	ReturnType string
	Name       string
	Params     []*Param
	Body       *SubroutineBody
}

type Param struct {
	Type string
	Name string
}

type SubroutineBody struct {
	Vars       []*LocalVarDeclaration
	Statements []Statement
}

type LocalVarDeclaration struct {
	Type string
	Name string
}

type Statement interface {
	statement() // marker method
}

type DoStatement struct {
	Call Expression
}

func (s *DoStatement) statement() {}

type LetStatement struct {
	Name  string
	Index Expression
	Value Expression
}

func (s *LetStatement) statement() {}

type WhileStatement struct {
	Condition Expression
	Body      []Statement
}

func (s *WhileStatement) statement() {}

type ReturnStatement struct {
	Value Expression
}

func (s *ReturnStatement) statement() {}

type IfStatement struct {
	Condition Expression
	Then      []Statement
	Else      []Statement
}

func (s *IfStatement) statement() {}

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
