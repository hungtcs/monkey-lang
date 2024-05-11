package syntax

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// 它首先检查了当前正在查看的字符l.ch，根据具体的字符来返回对应的词法单元。
// 在返回词法单元之前，位于所输入字符串中的指针会前移，所以之后再次调用NextToken()时，l.ch字段就已经更新过了。
// 最后，名为newToken的小型函数可以帮助初始化这些词法单元。
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	var tok Token
	switch l.ch {
	case '=':
		switch l.peekChar() {
		case '=':
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch)}
		default:
			tok = newToken(ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(PLUS, l.ch)
	case '-':
		tok = newToken(MINUS, l.ch)
	case '*':
		tok = newToken(ASTERISK, l.ch)
	case '/':
		tok = newToken(SLASH, l.ch)
	case '!':
		switch l.peekChar() {
		case '=':
			ch := l.ch
			l.readChar()
			tok = Token{Type: NE, Literal: string(ch) + string(l.ch)}
		default:
			tok = newToken(BANG, l.ch)
		}
	case '<':
		tok = newToken(LT, l.ch)
	case '>':
		tok = newToken(GT, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case ';':
		tok = newToken(SEMICOLON, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '{':
		tok = newToken(LBRACE, l.ch)
	case '}':
		tok = newToken(RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = INT
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// peekChar()与readChar()非常类似，但这个函数不会前移l.position和l.readPosition。
// 它的目的只是窥视一下输入中的下一个字符，不会移动位于输入中的指针位置，
// 这样就能知道下一步在调用readChar()时会返回什么。
// 大多数词法分析器和语法分析器具有这样的“窥视”函数，且大部分情况是用来向前看一个字符的。
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// readChar的目的是读取input中的下一个字符，并前移其在input中的位置。
// 这个过程的第一件事就是检查是否已经到达input的末尾。
// 如果是，则将l.ch设置为0，这是NUL字符的ASCII编码，用来表示“尚未读取任何内容”或“文件结尾”。
// 如果还没有到达input的末尾，则将l.ch设置为下一个字符，即l.input[l.readPosition]指向的字符。
//
// 之后，将l.position更新为刚用过的l.readPosition，然后将l.readPosition加1。
// 这样一来，l.readPosition就始终指向下一个将读取的字符位置，而l.position始终指向刚刚读取的位置。
// 这个特性很快就会派上用场。
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readIdentifier()函数顾名思义，就是读入一个标识符并前移词法分析器的扫描位置，直到遇见非字母字符。
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func newToken(t TokenType, ch byte) Token {
	return Token{
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
