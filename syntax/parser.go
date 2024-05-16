package syntax

import (
	"fmt"
	"strconv"
)

// 运算符优先级 (precedence)
const (
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESS_GREATER // > or <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X or !X
	CALL         // myFunction(X)
	INDEX        // a[i]
)

// 运算符对应的优先级
var precedence = map[Token]int{
	EQ:       EQUALS,
	NE:       EQUALS,
	LT:       LESS_GREATER,
	GT:       LESS_GREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	STAR:     PRODUCT,
	LPAREN:   CALL,
	LBRACKET: INDEX,
}

// 用于实现普拉特语法分析器
type (
	prefixParseFn func() Expr
	infixParseFn  func(Expr) Expr
)

type Parser struct {
	l *Lexer

	// pos    Position
	curTok TokenValue

	prefixParseFns map[Token]prefixParseFn
	infixParseFns  map[Token]infixParseFn
}

func (p *Parser) registerPrefixFn(tokenType Token, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType Token, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t Token) {
	panic(fmt.Errorf(
		`no prefix parse function for "%s" found`,
		t,
	))
}

func (p *Parser) nextToken() Position {
	pos := p.curTok.pos
	p.curTok = p.l.NextToken()
	return pos
}

func (p *Parser) curTokenIs(t Token) bool {
	return p.curTok.Type == t
}

// 断言当前 token 是否为 t
func (p *Parser) expect(t Token) {
	if !p.curTokenIs(t) {
		panic(
			NewError(
				p.curTok.pos,
				fmt.Sprintf(
					`expected next token to be "%s", got "%s" instead`,
					t, p.curTok,
				),
			),
		)
	}
}

// 断言当前 token 是否为 t，如果是，则前移
func (p *Parser) consume(t Token) Position {
	if p.curTok.Type == t {
		return p.nextToken()
	}
	panic(
		NewError(
			p.curTok.pos,
			fmt.Sprintf(
				`expected next token to be "%s", got "%s" instead`,
				t, p.curTok,
			),
		),
	)
}

// func (p *Parser) peekPrecedence() int {
// 	if val, ok := precedence[p.peekTok.Type]; ok {
// 		return val
// 	}
// 	return LOWEST
// }

func (p *Parser) curPrecedence() int {
	if val, ok := precedence[p.curTok.Type]; ok {
		return val
	}
	return LOWEST
}

func (p *Parser) Parse() (_ *Program, err error) {
	defer p.l.recover(&err)

	program := &Program{}
	program.Stmts = make([]Stmt, 0)

	for p.curTok.Type != EOF {
		stmt := p.parseStmt()
		program.Stmts = append(program.Stmts, stmt)
	}

	return program, nil
}

func (p *Parser) parseStmt() Stmt {
	switch p.curTok.Type {
	case LET:
		return p.parseLetStmt()
	case RETURN:
		return p.parseReturnStmt()
	default:
		return p.parseExprStmt()
	}
}

func (p *Parser) parseLetStmt() *LetStmt {
	pos := p.nextToken() // 消耗 let token
	stmt := &LetStmt{
		Pos: pos,
	}

	stmt.Name = &Identifier{Value: p.curTok.Literal}
	pos = p.consume(IDENT)
	stmt.Name.Pos = pos

	p.consume(ASSIGN)
	stmt.Value = p.parseExpr(LOWEST)
	if p.curTokenIs(SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStmt() *ReturnStmt {
	pos := p.nextToken()
	stmt := &ReturnStmt{
		Pos: pos,
	}
	stmt.Value = p.parseExpr(LOWEST)

	// 分号是可选的
	if p.curTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExprStmt() *ExprStmt {
	stmt := &ExprStmt{}
	stmt.Expr = p.parseExpr(LOWEST)

	// 分号是可选的
	if p.curTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpr(precedence int) Expr {
	prefix := p.prefixParseFns[p.curTok.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curTok.Type)
		return nil
	}
	leftExp := prefix()

	// 不是分号，并且优先级小于下一个词法单元的优先级，则继续解析
	for !p.curTokenIs(SEMICOLON) && precedence < p.curPrecedence() {
		infix := p.infixParseFns[p.curTok.Type]
		if infix == nil {
			return leftExp
		}
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expr {
	value := p.curTok.Literal
	pos := p.nextToken()
	return &Identifier{Pos: pos, Value: value}
}

func (p *Parser) parseIntegerLiteral() Expr {
	raw := p.curTok.Literal
	pos := p.nextToken()
	expr := &IntegerLiteral{Raw: raw, Pos: pos}
	value, err := strconv.ParseInt(raw, 0, 64)
	if err != nil {
		panic(fmt.Errorf("could not parse %q as integer", p.curTok.Literal))
	}
	expr.Value = value
	return expr
}

func (p *Parser) parseStringLiteral() Expr {
	val := p.curTok.Literal
	pos := p.nextToken()
	return &StringLiteral{pos: pos, Value: val}
}

func (p *Parser) parseBoolean() Expr {
	raw := p.curTok.Literal
	val := p.curTokenIs(TRUE)
	pos := p.nextToken()
	return &Boolean{pos: pos, raw: raw, Value: val}
}

func (p *Parser) parsePrefixExpr() Expr {
	op := p.curTok.Type
	pos := p.nextToken()
	expr := &PrefixExpr{
		Op:  op,
		Pos: pos,
	}
	expr.Right = p.parseExpr(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpr(left Expr) Expr {
	expr := &InfixExpr{
		Op:   p.curTok.Type,
		Left: left,
	}
	precedence := p.curPrecedence()
	p.nextToken() // 消耗运算符
	expr.Right = p.parseExpr(precedence)

	return expr
}

func (p *Parser) parseGroupedExpr() Expr {
	p.nextToken() // 消耗左括号
	exp := p.parseExpr(LOWEST)
	p.consume(RPAREN)
	return exp
}

func (p *Parser) parseArrayLiteral() Expr {
	start := p.nextToken()
	expr := &ArrayLiteral{start: start}
	expr.Items = p.parseExprList(RBRACKET)
	end := p.consume(RBRACKET)
	expr.end = end
	return expr
}

func (p *Parser) parseMapLiteral() Expr {
	start := p.nextToken()
	expr := &MapLiteral{start: start}
	expr.Pairs = make(map[Expr]Expr)

	// 检测到右括号，结束循环
	for !p.curTokenIs(RBRACE) {
		key := p.parseExpr(LOWEST) // 解析 Key
		p.consume(COLON)           // 解析冒号
		val := p.parseExpr(LOWEST) // 解析 Value
		expr.Pairs[key] = val

		// 如果下一个字符不是右括号，并且不是逗号，则结束循环
		if !p.curTokenIs(RBRACE) {
			p.consume(COMMA)
		}
	}

	end := p.consume(RBRACE)
	expr.end = end

	return expr
}

// 读取参数列表，知道遇到 end token，但是不消耗 end token
func (p *Parser) parseExprList(end Token) []Expr {
	exprs := make([]Expr, 0)
	if p.curTokenIs(end) {
		return exprs
	}
	exprs = append(exprs, p.parseExpr(LOWEST)) // 解析第一个参数
	for p.curTokenIs(COMMA) {
		p.consume(COMMA)
		exprs = append(exprs, p.parseExpr(LOWEST))
	}
	return exprs
}

func (p *Parser) parseBlockStmt() *BlockStmt {
	start := p.nextToken()
	block := &BlockStmt{start: start}
	block.Stmts = make([]Stmt, 0)
	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		stmt := p.parseStmt()
		block.Stmts = append(block.Stmts, stmt)
	}
	end := p.consume(RBRACE)
	block.end = end
	return block
}

func (p *Parser) parseIfExpr() Expr {
	pos := p.nextToken()
	expr := &IfExpr{pos: pos}
	p.consume(LPAREN)
	expr.Cond = p.parseExpr(LOWEST)
	p.consume(RPAREN)
	p.expect(LBRACE)
	expr.Consequence = p.parseBlockStmt()

	// 判断是否有 else 部分
	if p.curTokenIs(ELSE) {
		elsePos := p.nextToken()
		expr.elsePos = elsePos
		p.expect(LBRACE)
		expr.Alternative = p.parseBlockStmt()
	}

	return expr
}

func (p *Parser) parseFunctionParams() []*Identifier {
	identifiers := make([]*Identifier, 0)
	p.nextToken() // 消耗左括号
	// 如果是右括号，直接返回
	if p.curTokenIs(RPAREN) {
		p.consume(RPAREN)
		return identifiers
	}
	val := p.curTok.Literal
	pos := p.nextToken()
	identifier := &Identifier{Pos: pos, Value: val}
	identifiers = append(identifiers, identifier)
	for p.curTokenIs(COMMA) {
		p.consume(COMMA)
		val := p.curTok.Literal
		pos := p.nextToken()
		identifier := &Identifier{Pos: pos, Value: val}
		identifiers = append(identifiers, identifier)
	}
	p.consume(RPAREN)
	return identifiers
}

func (p *Parser) parseFunctionLiteral() Expr {
	pos := p.nextToken()
	expr := &FunctionLiteral{pos: pos}
	p.expect(LPAREN)
	expr.Params = p.parseFunctionParams()
	p.expect(LBRACE)
	expr.Body = p.parseBlockStmt()
	return expr
}

func (p *Parser) parseCallExpr(function Expr) Expr {
	start := p.consume(LPAREN)
	expr := &CallExpr{start: start, Function: function}
	expr.Args = p.parseExprList(RPAREN)
	end := p.consume(RPAREN)
	expr.end = end
	return expr
}

func (p *Parser) parseIndexExpr(left Expr) Expr {
	p.consume(LBRACKET)
	indexExpr := &IndexExpr{Left: left}
	indexExpr.Index = p.parseExpr(LOWEST)
	end := p.consume(RBRACKET)
	indexExpr.end = end
	return indexExpr
}

func NewParser(input string) *Parser {
	p := &Parser{
		l:              NewLexer(input),
		prefixParseFns: make(map[Token]prefixParseFn),
		infixParseFns:  make(map[Token]infixParseFn),
	}
	// 注册前缀解析函数
	p.registerPrefixFn(IDENT, p.parseIdentifier)
	p.registerPrefixFn(INT, p.parseIntegerLiteral)
	p.registerPrefixFn(TRUE, p.parseBoolean)
	p.registerPrefixFn(FALSE, p.parseBoolean)
	p.registerPrefixFn(BANG, p.parsePrefixExpr)
	p.registerPrefixFn(PLUS, p.parsePrefixExpr)
	p.registerPrefixFn(MINUS, p.parsePrefixExpr)
	p.registerPrefixFn(LPAREN, p.parseGroupedExpr)
	p.registerPrefixFn(LBRACE, p.parseMapLiteral)
	p.registerPrefixFn(LBRACKET, p.parseArrayLiteral)
	p.registerPrefixFn(IF, p.parseIfExpr)
	p.registerPrefixFn(FUNCTION, p.parseFunctionLiteral)
	p.registerPrefixFn(STRING, p.parseStringLiteral)

	// 注册中缀解析函数
	p.registerInfixFn(PLUS, p.parseInfixExpr)
	p.registerInfixFn(MINUS, p.parseInfixExpr)
	p.registerInfixFn(STAR, p.parseInfixExpr)
	p.registerInfixFn(SLASH, p.parseInfixExpr)
	p.registerInfixFn(EQ, p.parseInfixExpr)
	p.registerInfixFn(NE, p.parseInfixExpr)
	p.registerInfixFn(LT, p.parseInfixExpr)
	p.registerInfixFn(GT, p.parseInfixExpr)
	p.registerInfixFn(LPAREN, p.parseCallExpr)
	p.registerInfixFn(LBRACKET, p.parseIndexExpr)

	// read twice to set curTok and peekTok
	p.nextToken()

	return p
}
