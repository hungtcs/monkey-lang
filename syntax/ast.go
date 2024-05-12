package syntax

import (
	"bytes"
	"fmt"
	"strings"
)

type Node interface {
	fmt.Stringer
	Literal() string // 返回与其关联的词法单元字面量
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
	Tok   Token
	Value string
}

// String implements Expr.
func (i *Identifier) String() string {
	return i.Value
}

// Literal implements Expr.
func (i *Identifier) Literal() string {
	return i.Tok.Literal
}

// expr implements Expr.
func (i *Identifier) expr() {
	panic("unimplemented")
}

type LetStmt struct {
	Tok   Token
	Name  *Identifier
	Value Expr
}

// String implements Stmt.
func (l *LetStmt) String() string {
	var out bytes.Buffer
	out.WriteString(l.Tok.Literal + " ")
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
	return l.Tok.Literal
}

// stmt implements Stat.
func (l *LetStmt) stmt() {
	panic("unimplemented")
}

type ReturnStmt struct {
	Tok   Token
	Value Expr
}

// String implements Stmt.
func (r *ReturnStmt) String() string {
	var out bytes.Buffer
	out.WriteString(r.Tok.Literal + " ")
	if r.Value != nil {
		out.WriteString(r.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Literal implements Stmt.
func (r *ReturnStmt) Literal() string {
	return r.Tok.Literal
}

// stmt implements Stmt.
func (r *ReturnStmt) stmt() {
	panic("unimplemented")
}

type ExprStmt struct {
	Tok  Token
	Expr Expr
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
	return e.Tok.Literal
}

// stmt implements Stmt.
func (e *ExprStmt) stmt() {
	panic("unimplemented")
}

type IntegerLiteral struct {
	Tok   Token
	Value int64
}

// Literal implements Expr.
func (i *IntegerLiteral) Literal() string {
	return i.Tok.Literal
}

// String implements Expr.
func (i *IntegerLiteral) String() string {
	return i.Tok.Literal
}

// expr implements Expr.
func (i *IntegerLiteral) expr() {
	panic("unimplemented")
}

// 单目运算表达式
type PrefixExpr struct {
	Tok   Token
	Op    TokenType
	Right Expr
}

// Literal implements Expr.
func (p *PrefixExpr) Literal() string {
	return p.Tok.Literal
}

// String implements Expr.
func (p *PrefixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(string(p.Op))
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

// expr implements Expr.
func (p *PrefixExpr) expr() {
	panic("unimplemented")
}

type InfixExpr struct {
	Tok   Token
	Left  Expr
	Op    TokenType
	Right Expr
}

// Literal implements Expr.
func (i *InfixExpr) Literal() string {
	return i.Tok.Literal
}

// String implements Expr.
func (i *InfixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + string(i.Op) + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

// expr implements Expr.
func (i *InfixExpr) expr() {
	panic("unimplemented")
}

type Boolean struct {
	Tok   Token
	Value bool
}

// Literal implements Expr.
func (b *Boolean) Literal() string {
	return b.Tok.Literal
}

// String implements Expr.
func (b *Boolean) String() string {
	return b.Tok.Literal
}

// expr implements Expr.
func (b *Boolean) expr() {
	panic("unimplemented")
}

type BlockStmt struct {
	Tok   Token
	Stmts []Stmt
}

// Literal implements Stmt.
func (b *BlockStmt) Literal() string {
	return b.Tok.Literal
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
	Tok         Token
	Cond        Expr // 条件表达式
	Consequence *BlockStmt
	Alternative *BlockStmt
}

// Literal implements Expr.
func (i *IfExpr) Literal() string {
	return i.Tok.Literal
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
	Tok    Token
	Params []*Identifier
	Body   *BlockStmt
}

// Literal implements Expr.
func (f *FunctionLiteral) Literal() string {
	return f.Tok.Literal
}

// String implements Expr.
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}
	out.WriteString(f.Tok.Literal)
	out.WriteString("(")
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
	Tok      Token
	Function Expr
	Args     []Expr
}

// Literal implements Expr.
func (c *CallExpr) Literal() string {
	return c.Tok.Literal
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

var (
	_ Node = (*Program)(nil)
	_ Expr = (*Identifier)(nil)
	_ Stmt = (*LetStmt)(nil)
	_ Stmt = (*ReturnStmt)(nil)
	_ Stmt = (*ExprStmt)(nil)
	_ Expr = (*IntegerLiteral)(nil)
	_ Expr = (*PrefixExpr)(nil)
	_ Expr = (*InfixExpr)(nil)
	_ Expr = (*Boolean)(nil)
	_ Stmt = (*BlockStmt)(nil)
	_ Expr = (*IfExpr)(nil)
	_ Expr = (*FunctionLiteral)(nil)
	_ Expr = (*CallExpr)(nil)
)
