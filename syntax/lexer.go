package syntax

import (
	"fmt"
	"strings"
	"unicode"
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
	start := l.pos
	var tok TokenValue
	switch c {
	case '=', '!', '+', '-', '*', '/', '>', '<':
		l.nextChar()
		switch l.peekChar() {
		case '=':
			l.nextChar()
			switch c {
			case '=':
				tok = TokenValue{pos: start, Type: EQ, Literal: string(c) + "="}
			case '!':
				tok = TokenValue{pos: start, Type: NE, Literal: string(c) + "="}
			case '+':
				tok = TokenValue{pos: start, Type: PLUS, Literal: string(c) + "="}
			case '-':
				tok = TokenValue{pos: start, Type: MINUS, Literal: string(c) + "="}
			case '*':
				tok = TokenValue{pos: start, Type: STAR, Literal: string(c) + "="}
			case '/':
				tok = TokenValue{pos: start, Type: SLASH, Literal: string(c) + "="}
			case '>':
				tok = TokenValue{pos: start, Type: GE, Literal: string(c) + "="}
			case '<':
				tok = TokenValue{pos: start, Type: LE, Literal: string(c) + "="}
			}
		default:
			switch c {
			case '=':
				tok = createToken(ASSIGN, c, start)
			case '!':
				tok = createToken(BANG, c, start)
			case '+':
				tok = createToken(PLUS, c, start)
			case '-':
				tok = createToken(MINUS, c, start)
			case '*':
				tok = createToken(STAR, c, start)
			case '/':
				tok = createToken(SLASH, c, start)
			case '>':
				tok = createToken(GT, c, start)
			case '<':
				tok = createToken(LT, c, start)
			}
		}
	case ',':
		l.nextChar()
		tok = createToken(COMMA, c, start)
	case ':':
		l.nextChar()
		tok = createToken(COLON, c, start)
	case ';':
		l.nextChar()
		tok = createToken(SEMICOLON, c, start)
	case '(':
		l.nextChar()
		tok = createToken(LPAREN, c, start)
	case ')':
		l.nextChar()
		tok = createToken(RPAREN, c, start)
	case '{':
		l.nextChar()
		tok = createToken(LBRACE, c, start)
	case '}':
		l.nextChar()
		tok = createToken(RBRACE, c, start)
	case '[':
		l.nextChar()
		tok = createToken(LBRACKET, c, start)
	case ']':
		l.nextChar()
		tok = createToken(RBRACKET, c, start)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isIdentifierStart(c) {
			tok.pos = start
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
		} else if isDigit(c) {
			tok.pos = start
			tok.Literal = l.readNumber()
			tok.Type = INT
		} else {
			tok = createToken(ILLEGAL, c, start)
		}
	}
	tok.pos = start
	return tok
}

// 读取当前游标指向的字符
func (l *Lexer) peekChar() rune {
	if l.cursor >= len(l.input) {
		return 0
	}

	// ASCII
	if b := l.input[l.cursor]; b < utf8.RuneSelf {
		return rune(b)
	}

	r, _ := utf8.DecodeRuneInString(l.input[l.cursor:])
	return r
}

// 向后移动游标，并返回游标对应的字符
func (l *Lexer) nextChar() rune {
	// 到达文件末尾
	if l.cursor >= len(l.input) {
		return 0
	}

	var r rune
	var s int
	if b := l.input[l.cursor]; b < utf8.RuneSelf {
		s = 1
		r = rune(b)
	} else {
		r, s = utf8.DecodeRuneInString(l.input[l.cursor:])
	}

	l.cursor += s

	// 设置位置
	if r == '\n' {
		l.pos.Col = 1
		l.pos.Line += 1
	} else {
		l.pos.Col += 1
	}

	return r
}

func (l *Lexer) readNumber() string {
	position := l.cursor
	for c := l.peekChar(); isDigit(c); c = l.peekChar() {
		l.nextChar()
	}
	return l.input[position:l.cursor]
}

func (l *Lexer) readString() string {
	l.nextChar() // 消耗引号
	start := l.cursor
	for c := l.peekChar(); c != '"' && c != 0; c = l.peekChar() {
		l.nextChar()
	}
	val := l.input[start:l.cursor]
	l.nextChar() // 消耗引号
	return val
}

// readIdentifier()函数顾名思义，就是读入一个标识符并前移词法分析器的扫描位置，直到遇见非字母字符。
func (l *Lexer) readIdentifier() string {
	position := l.cursor
	for c := l.peekChar(); isIdentifier(c); c = l.peekChar() {
		l.nextChar()
	}
	return l.input[position:l.cursor]
}

func (l *Lexer) skipWhitespace() {
	for c := l.peekChar(); c == ' ' || c == '\t' || c == '\n' || c == '\r'; c = l.peekChar() {
		l.nextChar()
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

func createToken(t Token, ch rune, pos Position) TokenValue {
	return TokenValue{
		pos:     pos,
		Type:    t,
		Literal: string(ch),
	}
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isIdentifier(c rune) bool {
	return isIdentifierStart(c) || isDigit(c)
}

// 判断 c 是否为一个合法标识符的开始
func isIdentifierStart(c rune) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_' || unicode.IsLetter(c)
}
