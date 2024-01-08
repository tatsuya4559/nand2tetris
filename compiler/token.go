package main

type TokenKind int

const (
	TokenEOF TokenKind = iota

	TokenIdentifier
	TokenNumber
	TokenString

	/*** symbols ***/
	TokenLBrace
	TokenRBrace

	TokenLParen
	TokenRParen

	TokenLBracket
	TokenRBracket

	TokenDot
	TokenComma

	TokenSemicolon

	TokenPlus
	TokenMinus
	TokenAsterisk
	TokenSlash

	TokenAmpersand
	TokenVerticalLine
	TokenTilda

	TokenLT
	TokenGT
	TokenEqual

	/*** keywords ***/
	TokenClass
	TokenThis

	TokenMethod
	TokenFunction
	TokenConstructor

	TokenInt
	TokenBoolean
	TokenChar
	TokenVoid

	TokenVar
	TokenStatic
	TokenField

	TokenLet
	TokenDo
	TokenIf
	TokenElse
	TokenWhile
	TokenReturn

	TokenTrue
	TokenFalse
	TokenNull
)

type Token struct {
	Kind    TokenKind
	Literal string
}

var keywords = map[string]TokenKind{
	"class":       TokenClass,
	"this":        TokenThis,
	"method":      TokenMethod,
	"function":    TokenFunction,
	"constractor": TokenConstructor,
	"int":         TokenInt,
	"boolean":     TokenBoolean,
	"char":        TokenChar,
	"void":        TokenVoid,
	"var":         TokenVar,
	"static":      TokenStatic,
	"field":       TokenField,
	"let":         TokenLet,
	"do":          TokenDo,
	"if":          TokenIf,
	"else":        TokenElse,
	"while":       TokenWhile,
	"return":      TokenReturn,
	"true":        TokenTrue,
	"false":       TokenFalse,
	"null":        TokenNull,
}

func LookupKeyword(ident string) TokenKind {
	kw, isKw := keywords[ident]
	if isKw {
		return kw
	}
	return TokenIdentifier
}
