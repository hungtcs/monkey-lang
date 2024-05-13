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
)

// 运算符对应的优先级
var precedence = map[TokenType]int{
	EQ:       EQUALS,
	NE:       EQUALS,
	LT:       LESS_GREATER,
	GT:       LESS_GREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	ASTERISK: PRODUCT,
	LPAREN:   CALL,
}

// 用于实现普拉特语法分析器
type (
	prefixParseFn func() Expr
	infixParseFn  func(Expr) Expr
)

type Parser struct {
	l      *Lexer
	errors []string

	curTok  Token
	peekTok Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func (p *Parser) registerPrefixFn(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekError(t TokenType) {
	p.errors = append(
		p.errors,
		fmt.Sprintf(
			`expected next token to be "%s", got "%s" instead`,
			t, p.peekTok.Type,
		),
	)
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	p.errors = append(
		p.errors,
		fmt.Sprintf(
			`no prefix parse function for "%s" found`,
			t,
		),
	)
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curTok.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekTok.Type == t
}

// peekTokenIs 检查下一个词法单元是否为 t，
// 如果是，则前移词法单元，
// 否则返回 false
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if val, ok := precedence[p.peekTok.Type]; ok {
		return val
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if val, ok := precedence[p.curTok.Type]; ok {
		return val
	}
	return LOWEST
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Parse() *Program {
	program := &Program{}
	program.Stmts = make([]Stmt, 0)

	for p.curTok.Type != EOF {
		stmt := p.parseStmt()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
		p.nextToken()
	}

	return program
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
	stmt := &LetStmt{
		Tok: p.curTok,
	}
	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.Name = &Identifier{Tok: p.curTok, Value: p.curTok.Literal}

	if !p.expectPeek(ASSIGN) {
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpr(LOWEST)
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStmt() *ReturnStmt {
	stmt := &ReturnStmt{
		Tok: p.curTok,
	}
	// 消耗掉 return
	p.nextToken()
	stmt.Value = p.parseExpr(LOWEST)

	// 分号是可选的
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExprStmt() *ExprStmt {
	stmt := &ExprStmt{Tok: p.curTok}
	stmt.Expr = p.parseExpr(LOWEST)
	// 分号是可选的
	if p.peekTokenIs(SEMICOLON) {
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
	for !p.peekTokenIs(SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekTok.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expr {
	return &Identifier{Tok: p.curTok, Value: p.curTok.Literal}
}

func (p *Parser) parseIntegerLiteral() Expr {
	expr := &IntegerLiteral{Tok: p.curTok}
	value, err := strconv.ParseInt(p.curTok.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curTok.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	expr.Value = value
	return expr
}

func (p *Parser) parseStringLiteral() Expr {
	return &StringLiteral{Tok: p.curTok, Value: p.curTok.Literal}
}

func (p *Parser) parseBoolean() Expr {
	return &Boolean{Tok: p.curTok, Value: p.curTokenIs(TRUE)}
}

func (p *Parser) parsePrefixExpr() Expr {
	expr := &PrefixExpr{
		Tok: p.curTok,
		Op:  p.curTok.Type,
	}
	p.nextToken()
	expr.Right = p.parseExpr(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpr(left Expr) Expr {
	expr := &InfixExpr{
		Tok:  p.curTok,
		Op:   p.curTok.Type,
		Left: left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpr(precedence)
	return expr
}

func (p *Parser) parseGroupedExpr() Expr {
	p.nextToken()
	exp := p.parseExpr(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseBlockStmt() *BlockStmt {
	block := &BlockStmt{Tok: p.curTok}
	block.Stmts = make([]Stmt, 0)
	p.nextToken()
	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		stmt := p.parseStmt()
		block.Stmts = append(block.Stmts, stmt)
		// if stmt != nil {
		// 	block.Stmts = append(block.Stmts, stmt)
		// }
		p.nextToken()
	}
	return block
}

func (p *Parser) parseIfExpr() Expr {
	expr := &IfExpr{Tok: p.curTok}
	if !p.expectPeek(LPAREN) {
		return nil
	}
	p.nextToken() // 消耗左括号
	expr.Cond = p.parseExpr(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}
	if !p.expectPeek(LBRACE) {
		return nil
	}
	expr.Consequence = p.parseBlockStmt()

	// 判断是否有 else 部分
	if p.peekTokenIs(ELSE) {
		p.nextToken()
		if !p.expectPeek(LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStmt()
	}

	return expr
}

func (p *Parser) parseFunctionParams() []*Identifier {
	identifiers := make([]*Identifier, 0)
	// 如果是右括号，直接返回
	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &Identifier{Tok: p.curTok, Value: p.curTok.Literal}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Tok: p.curTok, Value: p.curTok.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseFunctionLiteral() Expr {
	expr := &FunctionLiteral{Tok: p.curTok}
	if !p.expectPeek(LPAREN) {
		return nil
	}
	expr.Params = p.parseFunctionParams()
	if !p.expectPeek(LBRACE) {
		return nil
	}
	expr.Body = p.parseBlockStmt()
	return expr
}

func (p *Parser) parseCallArgs() []Expr {
	args := make([]Expr, 0)
	// 如果是右括号，直接返回
	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken() // 消耗左括号
	args = append(args, p.parseExpr(LOWEST))
	// 如果是逗号，继续解析
	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpr(LOWEST))
	}
	if !p.expectPeek(RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseCallExpr(function Expr) Expr {
	expr := &CallExpr{Tok: p.curTok, Function: function}
	expr.Args = p.parseCallArgs()
	return expr
}

func NewParser(input string) *Parser {
	p := &Parser{
		l:              NewLexer(input),
		errors:         make([]string, 0),
		prefixParseFns: make(map[TokenType]prefixParseFn),
		infixParseFns:  make(map[TokenType]infixParseFn),
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
	p.registerPrefixFn(IF, p.parseIfExpr)
	p.registerPrefixFn(FUNCTION, p.parseFunctionLiteral)
	p.registerPrefixFn(STRING, p.parseStringLiteral)

	// 注册中缀解析函数
	p.registerInfixFn(PLUS, p.parseInfixExpr)
	p.registerInfixFn(MINUS, p.parseInfixExpr)
	p.registerInfixFn(ASTERISK, p.parseInfixExpr)
	p.registerInfixFn(SLASH, p.parseInfixExpr)
	p.registerInfixFn(EQ, p.parseInfixExpr)
	p.registerInfixFn(NE, p.parseInfixExpr)
	p.registerInfixFn(LT, p.parseInfixExpr)
	p.registerInfixFn(GT, p.parseInfixExpr)
	p.registerInfixFn(LPAREN, p.parseCallExpr)

	// read twice to set curTok and peekTok
	p.nextToken()
	p.nextToken()

	return p
}
