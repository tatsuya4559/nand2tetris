package main

import (
	"fmt"
	"io"
	"slices"
	"strconv"
)

type CompilationEngine struct {
	tokenizer *Tokenizer

	currentToken Token
	peekToken    Token
}

func NewCompilationEngine(input io.Reader) *CompilationEngine {
	engine := CompilationEngine{
		tokenizer: NewTokenizer(input),
	}
	engine.nextToken()
	engine.nextToken()
	return &engine
}

func (e *CompilationEngine) SetInput(input io.Reader) {
	e.tokenizer = NewTokenizer(input)
}

func (e *CompilationEngine) expectPeek(kinds ...TokenKind) {
	if slices.Contains(kinds, e.peekToken.Kind) {
		e.nextToken()
		return
	}
	// Die("Line %d: Unexpected token found. Expected %v, but got %v",
	// 	e.tokenizer.CurrentLineNum(), kinds, e.peekToken.Kind)
	panic(fmt.Sprintf("Line %d: Unexpected token found. Expected %v, but got %v",
		e.tokenizer.CurrentLineNum(), kinds, e.peekToken.Kind))
}

func (e *CompilationEngine) nextToken() {
	e.currentToken = e.peekToken
	e.peekToken = e.tokenizer.NextToken()
}

func (e *CompilationEngine) Compile() *ClassDeclaration {
	return e.compileClass()
}

// It is assumed currentToken is TokenClass when this method is called.
func (e *CompilationEngine) compileClass() *ClassDeclaration {
	Assert(e.currentToken.Kind == TokenClass)

	var class ClassDeclaration
	e.expectPeek(TokenIdentifier)
	class.Name = e.currentToken.Literal

	e.expectPeek(TokenLBrace)
	for e.peekToken.Kind == TokenStatic || e.peekToken.Kind == TokenField {
		e.nextToken() // read static or field
		dec := e.compileClassVarDec()
		class.Vars = append(class.Vars, dec)
	}
	for e.peekToken.Kind == TokenConstructor ||
		e.peekToken.Kind == TokenFunction ||
		e.peekToken.Kind == TokenMethod {
		e.nextToken() // read constructor or function or method
		sub := e.compileSubroutine()
		class.Subroutines = append(class.Subroutines, sub)
	}
	e.expectPeek(TokenRBrace)

	return &class
}

func (e *CompilationEngine) compileClassVarDec() *ClassVarDeclaration {
	var dec ClassVarDeclaration
	dec.StorageClass = e.currentToken.Literal // static or field

	e.expectPeek(TokenIdentifier, TokenInt, TokenChar, TokenBoolean)
	dec.Type = e.currentToken.Literal

	e.expectPeek(TokenIdentifier)
	dec.Names = append(dec.Names, e.currentToken.Literal)

	for e.peekToken.Kind == TokenComma {
		e.nextToken() // read comma
		e.expectPeek(TokenIdentifier)
		dec.Names = append(dec.Names, e.currentToken.Literal)
	}

	e.expectPeek(TokenSemicolon)
	return &dec
}

func (e *CompilationEngine) compileSubroutine() *SubroutineDeclaration {
	var dec SubroutineDeclaration
	dec.Kind = e.currentToken.Literal

	e.expectPeek(TokenVoid, TokenIdentifier, TokenInt, TokenChar, TokenBoolean)
	dec.ReturnType = e.currentToken.Literal

	e.expectPeek(TokenIdentifier)
	dec.Name = e.currentToken.Literal

	e.expectPeek(TokenLParen)
	dec.Params = e.compileParameterList()
	e.expectPeek(TokenLBrace)
	dec.Body = e.compileSubroutineBody()

	return &dec
}

func (e *CompilationEngine) compileParameterList() []*Param {
	Assert(e.currentToken.Kind == TokenLParen)

	var params []*Param
	for e.peekToken.Kind != TokenRParen {
		var param Param
		e.expectPeek(TokenIdentifier, TokenInt, TokenChar, TokenBoolean)
		param.Type = e.currentToken.Literal
		e.expectPeek(TokenIdentifier)
		param.Name = e.currentToken.Literal
		params = append(params, &param)
	}
	e.expectPeek(TokenRParen)
	return params
}

func (e *CompilationEngine) compileSubroutineBody() *SubroutineBody {
	var body SubroutineBody
	for e.peekToken.Kind == TokenVar {
		body.Vars = append(body.Vars, e.compileLocalVarDec())
	}
	for e.peekToken.Kind != TokenRBrace {
		body.Statements = append(body.Statements, e.compileStatement())
	}

	e.expectPeek(TokenRBrace)
	return &body
}

func (e *CompilationEngine) compileLocalVarDec() *LocalVarDeclaration {
	var dec LocalVarDeclaration
	e.expectPeek(TokenVar)

	e.expectPeek(TokenIdentifier, TokenInt, TokenChar, TokenBoolean)
	dec.Type = e.currentToken.Literal

	e.expectPeek(TokenIdentifier)
	dec.Names = append(dec.Names, e.currentToken.Literal)
	for e.peekToken.Kind != TokenSemicolon {
		e.expectPeek(TokenComma)
		e.expectPeek(TokenIdentifier)
		dec.Names = append(dec.Names, e.currentToken.Literal)
	}
	e.expectPeek(TokenSemicolon)

	return &dec
}

func (e *CompilationEngine) compileStatement() Statement {
	switch e.peekToken.Kind {
	case TokenDo:
		e.nextToken()
		return e.compileDo()
	case TokenLet:
		e.nextToken()
		return e.compileLet()
	case TokenWhile:
		e.nextToken()
		return e.compileWhile()
	case TokenReturn:
		e.nextToken()
		return e.compileReturn()
	case TokenIf:
		e.nextToken()
		return e.compileIf()
	default:
		Assert(false)
		return nil // unreachable
	}
}

func (e *CompilationEngine) compileDo() *DoStatement {
	var stmt DoStatement
	stmt.Call = e.compileExpression()
	if _, ok := stmt.Call.(*SubroutineCallExpression); !ok {
		Die("expected subroutine call")
	}
	e.expectPeek(TokenSemicolon)
	return &stmt
}

func (e *CompilationEngine) compileLet() *LetStatement {
	var stmt LetStatement
	e.expectPeek(TokenIdentifier)
	stmt.Name = e.currentToken.Literal

	if e.peekToken.Kind == TokenLBracket {
		stmt.Index = e.compileExpression()
	}

	e.expectPeek(TokenEqual)

	stmt.Value = e.compileExpression()

	e.expectPeek(TokenSemicolon)

	return &stmt
}

func (e *CompilationEngine) compileWhile() *WhileStatement {
	var stmt WhileStatement
	e.expectPeek(TokenLParen)
	stmt.Condition = e.compileExpression()
	e.expectPeek(TokenRParen)
	e.expectPeek(TokenLBrace)
	for e.peekToken.Kind != TokenRBrace {
		stmt.Body = append(stmt.Body, e.compileStatement())
	}
	e.expectPeek(TokenRBrace)
	return &stmt
}

func (e *CompilationEngine) compileReturn() *ReturnStatement {
	var stmt ReturnStatement
	if e.peekToken.Kind != TokenSemicolon {
		stmt.Value = e.compileExpression()
	}
	e.expectPeek(TokenSemicolon)
	return &stmt
}

func (e *CompilationEngine) compileIf() *IfStatement {
	var stmt IfStatement
	e.expectPeek(TokenLParen)
	stmt.Condition = e.compileExpression()
	e.expectPeek(TokenRParen)
	e.expectPeek(TokenLBrace)
	for e.peekToken.Kind != TokenRBrace {
		stmt.Then = append(stmt.Then, e.compileStatement())
	}
	e.expectPeek(TokenRBrace)
	if e.peekToken.Kind == TokenElse {
		e.expectPeek(TokenElse)
		e.expectPeek(TokenLBrace)
		for e.peekToken.Kind != TokenRBrace {
			stmt.Else = append(stmt.Else, e.compileStatement())
		}
		e.expectPeek(TokenRBrace)
	}
	return &stmt
}

func (e *CompilationEngine) compileExpression() Expression {
	var left Expression
	switch e.peekToken.Kind {
	case TokenNull:
		e.nextToken()
		left = NULL
	case TokenThis:
		e.nextToken()
		left = THIS
	case TokenTrue:
		e.nextToken()
		left = TRUE
	case TokenFalse:
		e.nextToken()
		left = FALSE
	case TokenIdentifier:
		e.nextToken()
		left = &IdentExpression{Value: e.currentToken.Literal}
	case TokenString:
		e.nextToken()
		left = &StringLiteralExpression{Value: e.currentToken.Literal}
	case TokenNumber:
		e.nextToken()
		v, err := strconv.Atoi(e.currentToken.Literal)
		if err != nil {
			panic("did not int")
		}
		left = &IntLiteralExpression{Value: v}
	case TokenLParen:
		e.nextToken()
		expr := e.compileExpression()
		e.expectPeek(TokenRParen)
		left = expr
	case TokenTilda, TokenMinus:
		e.nextToken()
		var expr PrefixExpression
		expr.Operator = e.currentToken.Literal
		expr.Right = e.compileExpression()
		left = &expr
	}

	if e.peekToken.Kind == TokenDot {
		e.nextToken()
		dot := &DotAccessExpression{Left: left}
		e.expectPeek(TokenIdentifier)
		dot.Right = &IdentExpression{Value: e.currentToken.Literal}
		left = dot
	}

	switch e.peekToken.Kind {
	case TokenPlus, TokenMinus, TokenAsterisk, TokenSlash,
		TokenAmpersand, TokenVerticalLine, TokenLT, TokenGT, TokenEqual:
		var inf InfixExpression
		inf.Left = left
		e.nextToken()
		inf.Operator = e.currentToken.Literal
		inf.Right = e.compileExpression()
		return &inf
	case TokenLParen:
		e.nextToken()
		var call SubroutineCallExpression
		call.Func = left
		for e.peekToken.Kind != TokenRParen {
			call.Args = append(call.Args, e.compileExpression())
			if e.peekToken.Kind == TokenComma {
				e.nextToken()
			} else {
				break
			}
		}
		e.expectPeek(TokenRParen)
		return &call
	}

	return left
}
