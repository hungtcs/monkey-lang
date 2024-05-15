package syntax

import (
	"bytes"
	"fmt"
	"strings"
)

type Node interface {
	fmt.Stringer
	Span() (start, end Position) // 返回节点的开始和结束位置
	Literal() string             // 返回与其关联的词法单元字面量
}

type Stmt interface {
	Node
	stmt() // 用于标记 Stmt 节点
}

type Expr interface {
	Node
	expr() // 用于标记 Expr 节点
}

type Program struct {
	Stmts []Stmt
}

// Span implements Node.
func (p *Program) Span() (start Position, end Position) {
	panic("unimplemented")
}

// String implements Node.
func (p *Program) String() string {
	var out bytes.Buffer
	for _, stmt := range p.Stmts {
		out.WriteString(stmt.String())
	}
	return out.String()
}

// Literal implements Node.
func (p *Program) Literal() string {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].Literal()
	}
	return ""
}

type Identifier struct {
	// Tok   TokenValue
	Pos   Position
	Value string
}

// Span implements Expr.
func (i *Identifier) Span() (start Position, end Position) {
	return i.Pos, i.Pos.add(i.Value)
}

// String implements Expr.
func (i *Identifier) String() string {
	return i.Value
}

// Literal implements Expr.
func (i *Identifier) Literal() string {
	return i.String()
}

// expr implements Expr.
func (i *Identifier) expr() {
	panic("unimplemented")
}

type LetStmt struct {
	// Tok   Token
	Pos   Position
	Name  *Identifier
	Value Expr
}

// Span implements Stmt.
func (l *LetStmt) Span() (start Position, end Position) {
	return l.Pos, l.Pos.add("let")
}

// String implements Stmt.
func (l *LetStmt) String() string {
	var out bytes.Buffer
	out.WriteString("let ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Literal implements Stat.
func (l *LetStmt) Literal() string {
	return "let"
}

// stmt implements Stat.
func (l *LetStmt) stmt() {
	panic("unimplemented")
}

type ReturnStmt struct {
	Pos   Position
	Value Expr
}

// Span implements Stmt.
func (r *ReturnStmt) Span() (start Position, end Position) {
	_, end = r.Value.Span()
	return r.Pos, end
}

// String implements Stmt.
func (r *ReturnStmt) String() string {
	var out bytes.Buffer
	out.WriteString("return ")
	if r.Value != nil {
		out.WriteString(r.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Literal implements Stmt.
func (r *ReturnStmt) Literal() string {
	return RETURN.String()
}

// stmt implements Stmt.
func (r *ReturnStmt) stmt() {
	panic("unimplemented")
}

type ExprStmt struct {
	// Tok TokenValue
	// Pos  Position
	Expr Expr
}

// Span implements Stmt.
func (e *ExprStmt) Span() (start Position, end Position) {
	return e.Expr.Span()
}

// String implements Stmt.
func (e *ExprStmt) String() string {
	if e.Expr != nil {
		return e.Expr.String()
	}
	return ""
}

// Literal implements Stmt.
func (e *ExprStmt) Literal() string {
	// return e.Tok.Literal
	return e.Expr.Literal()
}

// stmt implements Stmt.
func (e *ExprStmt) stmt() {
	panic("unimplemented")
}

type IntegerLiteral struct {
	Raw   string
	Pos   Position
	Value int64
}

// Span implements Expr.
func (i *IntegerLiteral) Span() (start Position, end Position) {
	return i.Pos, i.Pos.add(i.Raw)
}

// Literal implements Expr.
func (i *IntegerLiteral) Literal() string {
	return i.Raw
}

// String implements Expr.
func (i *IntegerLiteral) String() string {
	return i.Raw
}

// expr implements Expr.
func (i *IntegerLiteral) expr() {
	panic("unimplemented")
}

type StringLiteral struct {
	pos   Position
	Value string
}

// Span implements Expr.
func (s *StringLiteral) Span() (start Position, end Position) {
	return s.pos, s.pos.add(s.Value)
}

// Literal implements Expr.
func (s *StringLiteral) Literal() string {
	return s.Value
}

// String implements Expr.
func (s *StringLiteral) String() string {
	return s.Value
}

// expr implements Expr.
func (s *StringLiteral) expr() {
	panic("unimplemented")
}

// 单目运算表达式
type PrefixExpr struct {
	Op    Token
	Pos   Position
	Right Expr
}

// Span implements Expr.
func (p *PrefixExpr) Span() (start Position, end Position) {
	_, end = p.Right.Span()
	return p.Pos, end
}

// Literal implements Expr.
func (p *PrefixExpr) Literal() string {
	return p.String()
}

// String implements Expr.
func (p *PrefixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Op.String())
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

// expr implements Expr.
func (p *PrefixExpr) expr() {
	panic("unimplemented")
}

type InfixExpr struct {
	Left  Expr
	Op    Token
	Right Expr
}

// Span implements Expr.
func (i *InfixExpr) Span() (start Position, end Position) {
	start, _ = i.Left.Span()
	_, end = i.Right.Span()
	return start, end
}

// Literal implements Expr.
func (i *InfixExpr) Literal() string {
	return i.String()
}

// String implements Expr.
func (i *InfixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + string(i.Op.String()) + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

// expr implements Expr.
func (i *InfixExpr) expr() {
	panic("unimplemented")
}

type Boolean struct {
	pos   Position
	raw   string
	Value bool
}

// Span implements Expr.
func (b *Boolean) Span() (start Position, end Position) {
	return b.pos, b.pos.add(b.raw)
}

// Literal implements Expr.
func (b *Boolean) Literal() string {
	return b.raw
}

// String implements Expr.
func (b *Boolean) String() string {
	return b.raw
}

// expr implements Expr.
func (b *Boolean) expr() {
	panic("unimplemented")
}

type BlockStmt struct {
	start Position
	end   Position
	Stmts []Stmt
}

// Span implements Stmt.
func (b *BlockStmt) Span() (start Position, end Position) {
	return start, end
}

// Literal implements Stmt.
func (b *BlockStmt) Literal() string {
	return "if"
}

// String implements Stmt.
func (b *BlockStmt) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, stmt := range b.Stmts {
		out.WriteString(stmt.String())
	}
	out.WriteString("}")
	return out.String()
}

// stmt implements Stmt.
func (b *BlockStmt) stmt() {
	panic("unimplemented")
}

type IfExpr struct {
	pos         Position
	elsePos     Position
	Cond        Expr // 条件表达式
	Consequence *BlockStmt
	Alternative *BlockStmt
}

// Span implements Expr.
func (i *IfExpr) Span() (start Position, end Position) {
	body := i.Alternative
	if body == nil {
		body = i.Consequence
	}
	_, end = body.Stmts[len(body.Stmts)-1].Span()
	return i.pos, end
}

// Literal implements Expr.
func (i *IfExpr) Literal() string {
	return "if"
}

// String implements Expr.
func (i *IfExpr) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(i.Cond.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(i.Alternative.String())
	}
	return out.String()
}

// expr implements Expr.
func (i *IfExpr) expr() {
	panic("unimplemented")
}

type FunctionLiteral struct {
	pos    Position
	Params []*Identifier
	Body   *BlockStmt
}

// Span implements Expr.
func (f *FunctionLiteral) Span() (start Position, end Position) {
	panic("unimplemented")
}

// Literal implements Expr.
func (f *FunctionLiteral) Literal() string {
	return "fn"
}

// String implements Expr.
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}

// expr implements Expr.
func (f *FunctionLiteral) expr() {
	panic("unimplemented")
}

type CallExpr struct {
	start    Position
	end      Position
	Function Expr
	Args     []Expr
}

// Span implements Expr.
func (c *CallExpr) Span() (start Position, end Position) {
	return start, end
}

// Literal implements Expr.
func (c *CallExpr) Literal() string {
	return "call"
}

// String implements Expr.
func (c *CallExpr) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range c.Args {
		args = append(args, a.String())
	}
	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// expr implements Expr.
func (c *CallExpr) expr() {
	panic("unimplemented")
}

type ArrayLiteral struct {
	start Position
	end   Position
	Items []Expr
}

// Span implements Expr.
func (a *ArrayLiteral) Span() (start Position, end Position) {
	return start, end
}

// Literal implements Expr.
func (a *ArrayLiteral) Literal() string {
	return "["
}

// String implements Expr.
func (a *ArrayLiteral) String() string {
	var items = make([]string, len(a.Items))
	for _, item := range a.Items {
		items = append(items, item.String())
	}
	var out bytes.Buffer
	// out.WriteString(a.Tok.Literal)
	out.WriteString("[")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("]")
	return out.String()
}

// expr implements Expr.
func (a *ArrayLiteral) expr() {
	panic("unimplemented")
}

type MapLiteral struct {
	start Position
	end   Position
	Pairs map[Expr]Expr
}

// Span implements Expr.
func (m *MapLiteral) Span() (start Position, end Position) {
	return start, end
}

// Literal implements Expr.
func (m *MapLiteral) Literal() string {
	return m.String()
}

// String implements Expr.
func (m *MapLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for k, v := range m.Pairs {
		out.WriteString(k.String())
		out.WriteString(":")
		out.WriteString(v.String())
		out.WriteString(",")
	}
	out.WriteString("}")
	return out.String()
}

// expr implements Expr.
func (m *MapLiteral) expr() {
	panic("unimplemented")
}

type IndexExpr struct {
	end   Position
	Left  Expr
	Index Expr
}

// Span implements Expr.
func (i *IndexExpr) Span() (start Position, end Position) {
	start, _ = i.Left.Span()
	return start, i.end
}

// Literal implements Expr.
func (i *IndexExpr) Literal() string {
	return ""
}

// String implements Expr.
func (i *IndexExpr) String() string {
	var out bytes.Buffer
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("]")
	return out.String()
}

// expr implements Expr.
func (i *IndexExpr) expr() {
	panic("unimplemented")
}

var (
	_ Node = (*Program)(nil)
	_ Expr = (*Identifier)(nil)
	_ Stmt = (*LetStmt)(nil)
	_ Stmt = (*ReturnStmt)(nil)
	_ Stmt = (*ExprStmt)(nil)
	_ Expr = (*IntegerLiteral)(nil)
	_ Expr = (*StringLiteral)(nil)
	_ Expr = (*PrefixExpr)(nil)
	_ Expr = (*InfixExpr)(nil)
	_ Expr = (*Boolean)(nil)
	_ Stmt = (*BlockStmt)(nil)
	_ Expr = (*IfExpr)(nil)
	_ Expr = (*FunctionLiteral)(nil)
	_ Expr = (*CallExpr)(nil)
	_ Expr = (*ArrayLiteral)(nil)
	_ Expr = (*MapLiteral)(nil)
	_ Expr = (*IndexExpr)(nil)
)
