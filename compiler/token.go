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

func (k TokenKind) String() string {
	switch k {
	case TokenEOF:
		return "TokenEOF"
	case TokenIdentifier:
		return "TokenIdentifier"
	case TokenNumber:
		return "TokenNumber"
	case TokenString:
		return "TokenString"
	case TokenLBrace:
		return "TokenLBrace"
	case TokenRBrace:
		return "TokenRBrace"
	case TokenLParen:
		return "TokenLParen"
	case TokenRParen:
		return "TokenRParen"
	case TokenLBracket:
		return "TokenLBracket"
	case TokenRBracket:
		return "TokenRBracket"
	case TokenDot:
		return "TokenDot"
	case TokenComma:
		return "TokenComma"
	case TokenSemicolon:
		return "TokenSemicolon"
	case TokenPlus:
		return "TokenPlus"
	case TokenMinus:
		return "TokenMinus"
	case TokenAsterisk:
		return "TokenAsterisk"
	case TokenSlash:
		return "TokenSlash"
	case TokenAmpersand:
		return "TokenAmpersand"
	case TokenVerticalLine:
		return "TokenVerticalLine"
	case TokenTilda:
		return "TokenTilda"
	case TokenLT:
		return "TokenLT"
	case TokenGT:
		return "TokenGT"
	case TokenEqual:
		return "TokenEqual"
	case TokenClass:
		return "TokenClass"
	case TokenThis:
		return "TokenThis"
	case TokenMethod:
		return "TokenMethod"
	case TokenFunction:
		return "TokenFunction"
	case TokenConstructor:
		return "TokenConstructor"
	case TokenInt:
		return "TokenInt"
	case TokenBoolean:
		return "TokenBoolean"
	case TokenChar:
		return "TokenChar"
	case TokenVoid:
		return "TokenVoid"
	case TokenVar:
		return "TokenVar"
	case TokenStatic:
		return "TokenStatic"
	case TokenField:
		return "TokenField"
	case TokenLet:
		return "TokenLet"
	case TokenDo:
		return "TokenDo"
	case TokenIf:
		return "TokenIf"
	case TokenElse:
		return "TokenElse"
	case TokenWhile:
		return "TokenWhile"
	case TokenReturn:
		return "TokenReturn"
	case TokenTrue:
		return "TokenTrue"
	case TokenFalse:
		return "TokenFalse"
	case TokenNull:
		return "TokenNull"
	default:
		panic("unreachable")
	}
}

type Token struct {
	Kind    TokenKind
	Literal string
}

var keywords = map[string]TokenKind{
	"class":       TokenClass,
	"this":        TokenThis,
	"method":      TokenMethod,
	"function":    TokenFunction,
	"constructor": TokenConstructor,
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
