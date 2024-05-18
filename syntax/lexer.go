package syntax

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type ReadLineFunc func() (string, error)

type Lexer struct {
	pos      Position // 当前读取的位置
	rest     string
	readline ReadLineFunc
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
	c := l.peekRune()
	start := l.pos
	var tok TokenValue
	switch c {
	case '=', '!', '+', '-', '*', '/', '>', '<':
		l.nextRune()
		switch l.peekRune() {
		case '=':
			l.nextRune()
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
		l.nextRune()
		tok = createToken(COMMA, c, start)
	case ':':
		l.nextRune()
		tok = createToken(COLON, c, start)
	case ';':
		l.nextRune()
		tok = createToken(SEMICOLON, c, start)
	case '(':
		l.nextRune()
		tok = createToken(LPAREN, c, start)
	case ')':
		l.nextRune()
		tok = createToken(RPAREN, c, start)
	case '{':
		l.nextRune()
		tok = createToken(LBRACE, c, start)
	case '}':
		l.nextRune()
		tok = createToken(RBRACE, c, start)
	case '[':
		l.nextRune()
		tok = createToken(LBRACKET, c, start)
	case ']':
		l.nextRune()
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

func (l *Lexer) eof() bool {
	return len(l.rest) == 0 && !l.readLine()
}

func (l *Lexer) readLine() bool {
	if l.readline != nil {
		var err error
		l.rest, err = l.readline()
		if err != nil {
			panic(err) // EOF or ErrInterrupt
		}
		return len(l.rest) > 0
	}
	return false
}

// 读取当前游标指向的字符
func (l *Lexer) peekRune() rune {
	if l.eof() {
		return 0
	}

	// ASCII
	if b := l.rest[0]; b < utf8.RuneSelf {
		return rune(b)
	}

	r, _ := utf8.DecodeRuneInString(l.rest[0:])
	return r
}

// 向后移动游标，并返回游标对应的字符
func (l *Lexer) nextRune() rune {
	if len(l.rest) == 0 {
		if !l.readLine() {
			panic("internal scanner error: readRune at EOF")
		}
		// Redundant, but eliminates the bounds-check below.
		if len(l.rest) == 0 {
			return 0
		}
	}

	var r rune
	var s int
	if b := l.rest[0]; b < utf8.RuneSelf {
		s = 1
		r = rune(b)
	} else {
		r, s = utf8.DecodeRuneInString(l.rest[0:])
	}

	l.rest = l.rest[s:]

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
	raw := new(strings.Builder)
	for c := l.peekRune(); isDigit(c); c = l.peekRune() {
		raw.WriteRune(c)
		l.nextRune()
	}
	return raw.String()
}

func (l *Lexer) readString() string {
	l.nextRune() // 消耗引号
	raw := new(strings.Builder)
	for c := l.peekRune(); c != '"' && c != 0; c = l.peekRune() {
		raw.WriteRune(c)
		l.nextRune()
	}
	l.nextRune() // 消耗引号
	return raw.String()
}

// readIdentifier()函数顾名思义，就是读入一个标识符并前移词法分析器的扫描位置，直到遇见非字母字符。
func (l *Lexer) readIdentifier() string {
	raw := new(strings.Builder)
	for c := l.peekRune(); isIdentifier(c); c = l.peekRune() {
		raw.WriteRune(c)
		l.nextRune()
	}
	return raw.String()
}

func (l *Lexer) skipWhitespace() {
	for c := l.peekRune(); c == ' ' || c == '\t' || c == '\n' || c == '\r'; c = l.peekRune() {
		l.nextRune()
	}
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		pos:  MakePosition(nil, 1, 1), // 初始位置，第一行第一个字符
		rest: input,
	}
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
