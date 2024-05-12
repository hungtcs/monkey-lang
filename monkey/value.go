package monkey

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hungtcs/monkey-lang/syntax"
)

type Value interface {
	fmt.Stringer
	Type() string
	Truth() bool
}

type HasUnary interface {
	Value
	Unary(op syntax.TokenType) (_ Value, err error)
}

type HasBinary interface {
	Value
	Binary(op syntax.TokenType, y Value) (_ Value, err error)
}

type Comparable interface {
	Value
	Compare(op syntax.TokenType, y Value) (_ Value, err error)
}

// 有序类型，如数字
type TotallyOrdered interface {
	Value
	Cmp(y Value) (_ int, err error)
}

type NullType int

// String implements Value.
func (n NullType) String() string {
	return "null"
}

// Truth implements Value.
func (n NullType) Truth() bool {
	return false
}

// Type implements Value.
func (n NullType) Type() string {
	return "null"
}

const Null NullType = 0

type Int int64

// Cmp implements TotallyOrdered.
func (i Int) Cmp(y Value) (_ int, err error) {
	yv, ok := y.(Int)
	if !ok {
		return 0, fmt.Errorf("invalid cmp operator: %s %s %s", i, syntax.EQ, y)
	}
	if i > yv {
		return 1, nil
	} else if i < yv {
		return -1, nil
	}
	return 0, nil
}

// Binary implements HasBinary.
func (i Int) Binary(op syntax.TokenType, y Value) (_ Value, err error) {
	var yv Int
	if y, ok := y.(Int); !ok {
		return nil, fmt.Errorf("invalid binary operator: %s %s %s", i, op, y)
	} else {
		yv = y
	}

	switch op {
	case syntax.PLUS:
		return i + yv, nil
	case syntax.MINUS:
		return i - yv, nil
	case syntax.ASTERISK:
		return i * yv, nil
	case syntax.SLASH:
		return i / yv, nil
	default:
		return nil, nil
	}
}

// Unary implements HasUnary.
func (i Int) Unary(op syntax.TokenType) (_ Value, err error) {
	switch op {
	case syntax.MINUS:
		return -i, nil
	case syntax.PLUS:
		return i, nil
	default:
		return nil, nil
	}
}

// String implements Value.
func (i Int) String() string {
	return fmt.Sprintf("%d", i)
}

// Truth implements Value.
func (i Int) Truth() bool {
	return i != 0
}

// Type implements Value.
func (i Int) Type() string {
	return "int"
}

type Bool bool

// Compare implements Comparable.
func (b Bool) Compare(op syntax.TokenType, y_ Value) (_ Value, err error) {
	y, ok := y_.(Bool)
	if !ok {
		return nil, fmt.Errorf("invalid cmp operator: %s %s %s", b, op, y_)
	}

	return Bool(threeway(op, b2i(bool(b))-b2i(bool(y)))), nil
}

const (
	True  Bool = true
	False Bool = false
)

// String implements Value.
func (b Bool) String() string {
	return fmt.Sprintf("%t", b)
}

// Truth implements Value.
func (b Bool) Truth() bool {
	return bool(b)
}

// Type implements Value.
func (b Bool) Type() string {
	return "bool"
}

type returnValue struct {
	Value Value
}

// String implements Value.
func (r *returnValue) String() string {
	return r.Value.String()
}

// Truth implements Value.
func (r *returnValue) Truth() bool {
	return r.Value.Truth()
}

// Type implements Value.
func (r *returnValue) Type() string {
	return r.Value.Type()
}

type Function struct {
	Params []*syntax.Identifier
	Body   *syntax.BlockStmt
	Env    *Env
}

// String implements Value.
func (f *Function) String() string {
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

// Truth implements Value.
func (f *Function) Truth() bool {
	return true
}

// Type implements Value.
func (f *Function) Type() string {
	return "function"
}

var (
	_ Value          = NullType(0)
	_ Value          = Int(0)
	_ HasUnary       = Int(0)
	_ HasBinary      = Int(0)
	_ TotallyOrdered = Int(0)
	_ Value          = Bool(false)
	_ Comparable     = Bool(false)
	_ Value          = (*returnValue)(nil)
	_ Value          = (*Function)(nil)
)
