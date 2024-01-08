package main

type TokenKind int

const (
	TokenEOF TokenKind = iota
	TokenKeyword
	TokenSymbol
	TokenIdentifier
	TokenInt
	TokenString
)

type Token struct {
	Kind    TokenKind
	Literal string
}

type Keyword int

const (
	// class
	KwClass Keyword = iota
	KwThis

	// subroutine
	KwMethod
	KwFunction
	KwConstructor

	// type
	KwInt
	KwBoolean
	KwChar
	KwVoid

	// variable
	KwVar
	KwStatic
	KwField

	// statement
	KwLet
	KwDo
	KwIf
	KwElse
	KwWhile
	KwReturn

	// constant
	KwTrue
	KwFalse
	KwNull
)

var keywords = map[string]Keyword{
	"class":       KwClass,
	"this":        KwThis,
	"method":      KwMethod,
	"function":    KwFunction,
	"constractor": KwConstructor,
	"int":         KwInt,
	"boolean":     KwBoolean,
	"char":        KwChar,
	"void":        KwVoid,
	"var":         KwVar,
	"static":      KwStatic,
	"field":       KwField,
	"let":         KwLet,
	"do":          KwDo,
	"if":          KwIf,
	"else":        KwElse,
	"while":       KwWhile,
	"return":      KwReturn,
	"true":        KwTrue,
	"false":       KwFalse,
	"null":        KwNull,
}

func LookupKeyword(ident string) (kw Keyword, found bool) {
	kw, found = keywords[ident]
	return
}

func IsKeyword(ident string) bool {
	_, isKw := keywords[ident]
	return isKw
}
