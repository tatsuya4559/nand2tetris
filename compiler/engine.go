package main

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type CompilationEngine struct {
	tokenizer *Tokenizer
	vmwriter  *VMWriter
	symtable  *SymbolTable

	currentToken Token
	peekToken    Token
}

func NewCompilationEngine(input io.Reader, out io.Writer) *CompilationEngine {
	engine := CompilationEngine{
		tokenizer: NewTokenizer(input),
		vmwriter:  NewVMWriter(out),
		symtable:  NewSymbolTable(),
	}
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
	Die("Line %d: Unexpected token found. Expected %v, but got %v(%v)",
		e.tokenizer.CurrentLineNum(), kinds, e.peekToken.Kind, e.peekToken.Literal)
}

func (e *CompilationEngine) nextToken() {
	e.currentToken = e.peekToken
	e.peekToken = e.tokenizer.NextToken()
}

func (e *CompilationEngine) Compile() {
	e.compileClass()
}

func (e *CompilationEngine) compileClass() {
	e.expectPeek(TokenClass)
	e.expectPeek(TokenIdentifier)
	className := e.currentToken.Literal
	e.expectPeek(TokenLBrace)

	var nFields int
	for e.peekToken.Kind == TokenStatic || e.peekToken.Kind == TokenField {
		nFields += e.compileClassVarDec()
	}

	for e.peekToken.Kind == TokenConstructor ||
		e.peekToken.Kind == TokenFunction ||
		e.peekToken.Kind == TokenMethod {
		e.compileSubroutine(className, nFields)
	}
	e.expectPeek(TokenRBrace)
}

func (e *CompilationEngine) compileClassVarDec() int {
	var nFields int
	e.expectPeek(TokenField, TokenStatic)
	storage := e.currentToken.Literal

	e.expectPeek(TokenIdentifier, TokenInt, TokenChar, TokenBoolean)
	typ := e.currentToken.Literal

	e.expectPeek(TokenIdentifier)
	name := e.currentToken.Literal

	e.symtable.Define(name, typ, storageToScope(storage))
	if storage == "field" {
		nFields++
	}

	for e.peekToken.Kind == TokenComma {
		e.expectPeek(TokenComma)
		e.expectPeek(TokenIdentifier)
		name = e.currentToken.Literal
		e.symtable.Define(name, typ, storageToScope(storage))
		nFields++
	}

	e.expectPeek(TokenSemicolon)
	return nFields
}

func (e *CompilationEngine) compileSubroutine(className string, nFields int) {
	e.symtable.ResetLocalScope()

	/* Declaration */
	e.expectPeek(TokenConstructor, TokenFunction, TokenMethod)
	kind := e.currentToken.Literal

	e.expectPeek(TokenVoid, TokenIdentifier, TokenInt, TokenChar, TokenBoolean)

	e.expectPeek(TokenIdentifier)
	name := fmt.Sprintf("%s.%s", className, e.currentToken.Literal)

	/* Params */
	e.compileParameterList()

	/* Body */
	e.expectPeek(TokenLBrace)

	nLocals := 0
	for e.peekToken.Kind == TokenVar {
		nLocals += e.compileLocalVarDec()
	}
	e.vmwriter.WriteFunction(name, nLocals)

	if kind == "constructor" {
		// this = Memory.alloc(nFields)
		e.vmwriter.WritePush(SegConst, nFields)
		e.vmwriter.WriteCall("Memory.alloc", 1)
		e.vmwriter.WritePop(SegPointer, 0)
	}

	for e.peekToken.Kind != TokenRBrace {
		e.compileStatement()
	}

	e.expectPeek(TokenRBrace)
}

func (e *CompilationEngine) compileParameterList() {
	e.expectPeek(TokenLParen)
	for e.peekToken.Kind != TokenRParen {
		e.expectPeek(TokenIdentifier, TokenInt, TokenChar, TokenBoolean)
		typ := e.currentToken.Literal

		e.expectPeek(TokenIdentifier)
		name := e.currentToken.Literal

		e.symtable.Define(name, typ, ScopeArg)

		if e.peekToken.Kind == TokenComma {
			e.expectPeek(TokenComma)
		} else {
			break
		}
	}
	e.expectPeek(TokenRParen)
}

func (e *CompilationEngine) compileLocalVarDec() int {
	var nLocals int
	e.expectPeek(TokenVar)

	e.expectPeek(TokenIdentifier, TokenInt, TokenChar, TokenBoolean)
	typ := e.currentToken.Literal

	e.expectPeek(TokenIdentifier)
	name := e.currentToken.Literal

	e.symtable.Define(name, typ, ScopeVar)
	nLocals++

	for e.peekToken.Kind != TokenSemicolon {
		e.expectPeek(TokenComma)
		e.expectPeek(TokenIdentifier)
		name = e.currentToken.Literal
		e.symtable.Define(name, typ, ScopeVar)
		nLocals++
	}
	e.expectPeek(TokenSemicolon)

	return nLocals
}

func (e *CompilationEngine) compileStatement() {
	switch e.peekToken.Kind {
	case TokenDo:
		e.compileDo()
	case TokenLet:
		e.compileLet()
	case TokenWhile:
		e.compileWhile()
	case TokenReturn:
		e.compileReturn()
	case TokenIf:
		e.compileIf()
	default:
		panic("unknown statement found")
	}
}

func (e *CompilationEngine) compileDo() {
	e.expectPeek(TokenDo)
	e.compileExpression()
	e.expectPeek(TokenSemicolon)
	e.vmwriter.WritePop(SegTemp, 0) // discard return value
}

func (e *CompilationEngine) compileLet() {
	var leftIsArray bool

	e.expectPeek(TokenLet)
	e.expectPeek(TokenIdentifier)
	name := e.currentToken.Literal
	entry := e.symtable.Find(name)

	if e.peekToken.Kind == TokenLBracket {
		leftIsArray = true
		e.expectPeek(TokenLBracket)

		e.vmwriter.WritePush(Segment(entry.Scope), entry.Index)
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdAdd)
		e.vmwriter.WritePop(SegPointer, 1)

		e.expectPeek(TokenRBracket)
	}

	e.expectPeek(TokenEqual)
	e.compileExpression()
	e.expectPeek(TokenSemicolon)

	if leftIsArray {
		e.vmwriter.WritePop(SegThat, 0)
	} else {
		e.vmwriter.WritePop(Segment(entry.Scope), entry.Index)
	}
}

func (e *CompilationEngine) compileWhile() {
	e.expectPeek(TokenWhile)

	loopLabel := genLabel("loop")
	endLabel := genLabel("end")

	/* Condition */
	e.vmwriter.WriteLabel(loopLabel)
	e.expectPeek(TokenLParen)
	e.compileExpression()
	e.vmwriter.WriteArithmeric(CmdNot)
	e.vmwriter.WriteIf(endLabel)
	e.expectPeek(TokenRParen)

	/* Body */
	e.expectPeek(TokenLBrace)
	for e.peekToken.Kind != TokenRBrace {
		e.compileStatement()
	}
	e.expectPeek(TokenRBrace)
	e.vmwriter.WriteGoto(loopLabel)
	e.vmwriter.WriteLabel(endLabel)
}

func (e *CompilationEngine) compileReturn() {
	e.expectPeek(TokenReturn)
	if e.peekToken.Kind != TokenSemicolon {
		e.compileExpression()
	} else {
		// pseudo return value for void function
		e.vmwriter.WritePush(SegConst, 0)
	}
	e.vmwriter.WriteReturn()
	e.expectPeek(TokenSemicolon)
}

func (e *CompilationEngine) compileIf() {
	e.expectPeek(TokenIf)

	e.expectPeek(TokenLParen)
	e.compileExpression()
	e.expectPeek(TokenRParen)

	elseLabel := genLabel("else")
	endLabel := genLabel("end")
	e.vmwriter.WriteArithmeric(CmdNot)
	e.vmwriter.WriteIf(elseLabel)

	e.expectPeek(TokenLBrace)
	for e.peekToken.Kind != TokenRBrace {
		e.compileStatement()
	}
	e.expectPeek(TokenRBrace)

	e.vmwriter.WriteGoto(endLabel)
	e.vmwriter.WriteLabel(elseLabel)
	if e.peekToken.Kind == TokenElse {
		e.expectPeek(TokenElse)
		e.expectPeek(TokenLBrace)
		for e.peekToken.Kind != TokenRBrace {
			e.compileStatement()
		}
		e.expectPeek(TokenRBrace)
	}
	e.vmwriter.WriteLabel(endLabel)
}

func (e *CompilationEngine) compileExpression() {
	var fn string
	var entry *SymbolTableEntry
	switch e.peekToken.Kind {
	case TokenNull:
		e.nextToken()
		e.vmwriter.WritePush(SegConst, 0)
	case TokenThis:
		e.nextToken()
		e.vmwriter.WritePush(SegPointer, 0)
	case TokenTrue:
		e.nextToken()
		e.vmwriter.WritePush(SegConst, 1)
		e.vmwriter.WriteArithmeric(CmdNeg)
	case TokenFalse:
		e.nextToken()
		e.vmwriter.WritePush(SegConst, 0)
	case TokenIdentifier:
		e.nextToken()
		entry = e.symtable.Find(e.currentToken.Literal)
		if entry == nil {
			fn = e.currentToken.Literal
			break
		}
		e.vmwriter.WritePush(Segment(entry.Scope), entry.Index)
	case TokenString:
		e.nextToken()
		str := e.currentToken.Literal
		e.vmwriter.WritePush(SegConst, len(str))
		e.vmwriter.WriteCall("String.new", 1)
		for _, r := range str {
			e.vmwriter.WritePush(SegConst, int(r))
			e.vmwriter.WriteCall("String.appendChar", 2)
		}
	case TokenNumber:
		e.nextToken()
		v, err := strconv.Atoi(e.currentToken.Literal)
		if err != nil {
			panic("did not int")
		}
		e.vmwriter.WritePush(SegConst, v)
	case TokenLParen:
		e.expectPeek(TokenLParen)
		e.compileExpression()
		e.expectPeek(TokenRParen)
	case TokenTilda:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdNot)
	case TokenMinus:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdNeg)
	default:
		Die("invalid expression: %v", e.peekToken)
	}

	if e.peekToken.Kind == TokenDot {
		e.expectPeek(TokenDot)
		e.expectPeek(TokenIdentifier)
		fn += "."
		fn += e.currentToken.Literal
	}
	if e.peekToken.Kind == TokenLBracket {
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdAdd)
		e.vmwriter.WritePop(SegPointer, 1)
		e.vmwriter.WritePush(SegThat, 0)
		e.expectPeek(TokenRBracket)
	}

	switch e.peekToken.Kind {
	case TokenPlus:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdAdd)
	case TokenMinus:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdSub)
	case TokenAsterisk:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteCall("Math.multiply", 2)
	case TokenSlash:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteCall("Math.divide", 2)
	case TokenAmpersand:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdAnd)
	case TokenVerticalLine:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdOr)
	case TokenLT:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdLt)
	case TokenGT:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdGt)
	case TokenEqual:
		e.nextToken()
		e.compileExpression()
		e.vmwriter.WriteArithmeric(CmdEq)
	case TokenLParen:
		e.expectPeek(TokenLParen)
		nArgs := 0
		if strings.HasPrefix(fn, ".") { /* when fn is a method */
			// append class name
			fn = entry.Type + fn
			// instance is already on the top of stack
			nArgs++
		}
		for e.peekToken.Kind != TokenRParen {
			e.compileExpression()
			nArgs++
			if e.peekToken.Kind == TokenComma {
				e.nextToken()
			} else {
				break
			}
		}
		e.expectPeek(TokenRParen)
		e.vmwriter.WriteCall(fn, nArgs)
	}
}

var labelSequence = 0

func genLabel(prefix string) string {
	label := fmt.Sprintf("%s.%d", prefix, labelSequence)
	labelSequence++
	return label
}

func storageToScope(storage string) Scope {
	switch storage {
	case "field":
		return ScopeField
	case "static":
		return ScopeStatic
	default:
		panic("unknown storage class")
	}
}
