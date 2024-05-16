package syntax

import "fmt"

type Token int8

func (t Token) String() string {
	return tokenNames[t]
}

const (
	ILLEGAL Token = iota
	EOF

	IDENT
	INT
	STRING

	ASSIGN // =
	PLUS   // +
	MINUS  // -
	STAR   // *
	SLASH  // /
	BANG   // ！
	LT     // <
	LE     // <=
	GT     // >
	GE     // >=
	EQ     // ==
	NE     // !=

	COLON     // :
	COMMA     // ,
	SEMICOLON // ;

	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	LET      // let
	IF       // if
	ELSE     // else
	TRUE     // true
	FALSE    // false
	RETURN   // return
	FUNCTION // fn
)

var tokenNames = [...]string{
	ILLEGAL: "illegal token",
	EOF:     "end of file",
	IDENT:   "identifier",
	INT:     "int",
	STRING:  "string",

	ASSIGN: "=",
	PLUS:   "+",
	MINUS:  "-",
	STAR:   "*",
	SLASH:  "/",
	BANG:   "!",
	LT:     "<",
	LE:     "<=",
	GT:     ">",
	GE:     ">=",
	EQ:     "==",
	NE:     "!=",

	COLON:     ":",
	COMMA:     ",",
	SEMICOLON: ";",

	LPAREN:   "(",
	RPAREN:   ")",
	LBRACE:   "{",
	RBRACE:   "}",
	LBRACKET: "[",
	RBRACKET: "]",

	LET:      "let",
	IF:       "if",
	ELSE:     "else",
	TRUE:     "true",
	FALSE:    "false",
	RETURN:   "return",
	FUNCTION: "fn",
}

var keywords = map[string]Token{
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
func LookupIdent(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type TokenValue struct {
	pos     Position
	Type    Token
	Literal string
}

func (t TokenValue) String() string {
	return fmt.Sprintf(`%s(literal="%s")`, t.Type, t.Literal)
}
