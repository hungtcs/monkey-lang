package syntax

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// 位置描述输入符文的位置。
type Position struct {
	file *string // filename
	Line int32   // 1-based line number; 0 if line unknown
	Col  int32   // 1-based column (rune) number; 0 if column unknown
}

func (p Position) add(s string) Position {
	if n := strings.Count(s, "\n"); n > 0 {
		p.Line += int32(n)
		s = s[strings.LastIndex(s, "\n")+1:]
		p.Col = 1
	}
	p.Col += int32(utf8.RuneCountInString(s))
	return p
}

func (p Position) Filename() string {
	if p.file != nil {
		return *p.file
	}
	return "<invalid>"
}

func (p Position) String() string {
	file := p.Filename()
	if p.Line > 0 {
		if p.Col > 0 {
			return fmt.Sprintf("%s:%d:%d", file, p.Line, p.Col)
		}
		return fmt.Sprintf("%s:%d", file, p.Line)
	}
	return file
}

// MakePosition returns position with the specified components.
func MakePosition(file *string, line, col int32) Position {
	return Position{file, line, col}
}

type Lexer struct {
	pos    Position // 当前读取的位置
	input  string
	cursor int
}

func (p *Lexer) recover(err *error) {
	switch e := recover().(type) {
	case nil:
	case error:
		*err = e
	default:
		*err = fmt.Errorf("parser panic: %v", e)
	}
}

// 它首先检查了当前正在查看的字符l.ch，根据具体的字符来返回对应的词法单元。
// 在返回词法单元之前，位于所输入字符串中的指针会前移，所以之后再次调用NextToken()时，l.ch字段就已经更新过了。
// 最后，名为newToken的小型函数可以帮助初始化这些词法单元。
func (l *Lexer) NextToken() TokenValue {
	l.skipWhitespace() // 跳过空白字符
	c := l.peekChar()

	var tok TokenValue
	switch c {
	case '=', '!', '+', '-', '*', '/', '>', '<':
		switch l.readChar() {
		case '=':
			l.readChar() // 消耗 ==
			switch c {
			case '=':
				tok = TokenValue{Type: EQ, Literal: string(c) + "="}
			case '!':
				tok = TokenValue{Type: NE, Literal: string(c) + "="}
			case '+':
				tok = TokenValue{Type: PLUS, Literal: string(c) + "="}
			case '-':
				tok = TokenValue{Type: MINUS, Literal: string(c) + "="}
			case '*':
				tok = TokenValue{Type: STAR, Literal: string(c) + "="}
			case '/':
				tok = TokenValue{Type: SLASH, Literal: string(c) + "="}
			case '>':
				tok = TokenValue{Type: GE, Literal: string(c) + "="}
			case '<':
				tok = TokenValue{Type: LE, Literal: string(c) + "="}
			}
		default:
			switch c {
			case '=':
				tok = newToken(ASSIGN, c)
			case '!':
				tok = newToken(BANG, c)
			case '+':
				tok = newToken(PLUS, c)
			case '-':
				tok = newToken(MINUS, c)
			case '*':
				tok = newToken(STAR, c)
			case '/':
				tok = newToken(SLASH, c)
			case '>':
				tok = newToken(GT, c)
			case '<':
				tok = newToken(LT, c)
			}
		}
	case ',':
		l.readChar()
		tok = newToken(COMMA, c)
	case ':':
		l.readChar()
		tok = newToken(COLON, c)
	case ';':
		l.readChar()
		tok = newToken(SEMICOLON, c)
	case '(':
		l.readChar()
		tok = newToken(LPAREN, c)
	case ')':
		l.readChar()
		tok = newToken(RPAREN, c)
	case '{':
		l.readChar()
		tok = newToken(LBRACE, c)
	case '}':
		l.readChar()
		tok = newToken(RBRACE, c)
	case '[':
		l.readChar()
		tok = newToken(LBRACKET, c)
	case ']':
		l.readChar()
		tok = newToken(RBRACKET, c)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(c) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(c) {
			tok.Literal = l.readNumber()
			tok.Type = INT
			return tok
		} else {
			tok = newToken(ILLEGAL, c)
		}
	}
	return tok
}

// 读取当前游标指向的字符
func (l *Lexer) peekChar() byte {
	if l.cursor >= len(l.input) {
		return 0
	} else {
		return l.input[l.cursor]
	}
}

// 向后移动游标，并返回游标对应的字符
func (l *Lexer) readChar() byte {
	// 到达文件末尾
	if l.cursor >= len(l.input) {
		return 0
	}

	l.cursor += 1
	c := l.peekChar()

	// 设置位置
	if c == '\n' {
		l.pos.Col = 1
		l.pos.Line += 1
	} else {
		l.pos.Col += 1
	}
	return c
}

func (l *Lexer) readNumber() string {
	position := l.cursor
	for c := l.peekChar(); isDigit(c); c = l.peekChar() {
		l.readChar()
	}
	return l.input[position:l.cursor]
}

func (l *Lexer) readString() string {
	l.readChar() // 消耗引号
	start := l.cursor
	for c := l.peekChar(); c != '"' && c != 0; c = l.peekChar() {
		l.readChar()
	}
	val := l.input[start:l.cursor]
	l.readChar() // 消耗引号
	return val
}

// readIdentifier()函数顾名思义，就是读入一个标识符并前移词法分析器的扫描位置，直到遇见非字母字符。
func (l *Lexer) readIdentifier() string {
	position := l.cursor
	for c := l.peekChar(); isLetter(c); c = l.peekChar() {
		l.readChar()
	}
	return l.input[position:l.cursor]
}

func (l *Lexer) skipWhitespace() {
	for c := l.peekChar(); c == ' ' || c == '\t' || c == '\n' || c == '\r'; c = l.peekChar() {
		l.readChar()
	}
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		pos:   MakePosition(nil, 1, 1), // 初始位置，第一行第一个字符
		input: input,
	}
	// 读入第一个字符
	// l.readChar()
	return l
}

func newToken(t Token, ch byte) TokenValue {
	return TokenValue{
		Type:    t,
		Literal: string(ch),
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isLetter辅助函数用来判断给定的参数是否为字母。
// 值得注意的是，这个函数虽然看起来简短，但意义重大，其决定了解释器所能处理的语言形式。
// 比如示例中包含ch =='_'，这意味着下划线_会被视为字母，允许在标识符和关键字中使用。
// 因此可以使用诸如foo_bar之类的变量名。
// 其他编程语言甚至允许在标识符中使用问号和感叹号。
// 如果读者也想这么做，那么可以修改这个isLetter函数。
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
