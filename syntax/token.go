package syntax

import "fmt"

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	STRING TokenType = "STRING"

	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	BANG     TokenType = "!"
	LT       TokenType = "<"
	LE       TokenType = "<="
	GT       TokenType = ">"
	GE       TokenType = ">="
	EQ       TokenType = "=="
	NE       TokenType = "!="

	COLON     TokenType = ":"
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	LET      TokenType = "LET"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	RETURN   TokenType = "RETURN"
	FUNCTION TokenType = "FUNCTION"
)

var keywords = map[string]TokenType{
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
	"fn":     FUNCTION,
}

// LookupIdent通过检查关键字表来判断给定的标识符是否是关键字。
// 如果是，则返回关键字的TokenType常量。
// 如果不是，则返回token.IDENT，这个TokenType表示当前是用户定义的标识符。
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf(`%s(literal="%s")`, t.Type, t.Literal)
}
