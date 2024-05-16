package monkey

import (
	"bytes"
	"fmt"
	"hash/maphash"
	"math/big"
	"strings"
	"unicode/utf8"

	"github.com/hungtcs/monkey-lang/syntax"
)

type Side bool

const (
	Left  Side = false
	Right Side = true
)

type Value interface {
	fmt.Stringer
	Type() string
	Hash() (uint32, error)
	Truth() bool
}

type Iterable interface {
	Value
}

type Sequence interface {
	Iterable
	Len() int
}

type Mapping interface {
	Value
	Get(v Value) (_ Value, _ bool, err error)
}

type Indexable interface {
	Value
	Len() int
	Index(i int) Value
}

type HasUnary interface {
	Value
	Unary(op syntax.Token) (_ Value, err error)
}

type HasBinary interface {
	Value
	Binary(op syntax.Token, y Value, side Side) (_ Value, err error)
}

type Comparable interface {
	Value
	Compare(op syntax.Token, y Value) (_ Value, err error)
}

// 有序类型，如数字
type TotallyOrdered interface {
	Value
	Cmp(y Value) (_ int, err error)
}

type Callable interface {
	Value
	Name() string
	CallInternal(args ...Value) (_ Value, err error)
}

type NullType int

// Hash implements Value.
func (n NullType) Hash() (uint32, error) {
	return 0, nil
}

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

// Hash implements Value.
func (i Int) Hash() (uint32, error) {
	lo := big.Word(i)
	return 12582917 * uint32(lo+3), nil
}

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
func (i Int) Binary(op syntax.Token, y Value, side Side) (_ Value, err error) {
	var yv Int
	if y, ok := y.(Int); !ok {
		return nil, nil
	} else {
		yv = y
	}

	switch op {
	case syntax.PLUS:
		return i + yv, nil
	case syntax.MINUS:
		return i - yv, nil
	case syntax.STAR:
		return i * yv, nil
	case syntax.SLASH:
		return i / yv, nil
	default:
		return nil, nil
	}
}

// Unary implements HasUnary.
func (i Int) Unary(op syntax.Token) (_ Value, err error) {
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

// Hash implements Value.
func (b Bool) Hash() (uint32, error) {
	return uint32(b2i(bool(b))), nil
}

// Compare implements Comparable.
func (b Bool) Compare(op syntax.Token, y_ Value) (_ Value, err error) {
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

type String string

var seed = maphash.MakeSeed()

// Hash implements Value.
func (s String) Hash() (uint32, error) {
	if len(s) >= 12 {
		// Call the Go runtime's optimized hash implementation,
		// which uses the AES instructions on amd64 and arm64 machines.
		h := maphash.String(seed, string(s))
		return uint32(h>>32) | uint32(h), nil
	}
	return softHashString(string(s)), nil
}

// softHashString computes the 32-bit FNV-1a hash of s in software.
func softHashString(s string) uint32 {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

// Index implements Indexable.
func (s String) Index(i int) Value {
	return s[i : i+1]
}

// Len implements Indexable.
func (s String) Len() int {
	return utf8.RuneCountInString(string(s))
}

// Binary implements HasBinary.
func (s String) Binary(op syntax.Token, y Value, side Side) (_ Value, err error) {
	switch op {
	case syntax.PLUS:
		switch y := y.(type) {
		case String:
			return s + y, nil
		}
	}
	return nil, nil
}

// String implements Value.
func (s String) String() string {
	return string(s)
}

// Truth implements Value.
func (s String) Truth() bool {
	return len(s) > 0
}

// Type implements Value.
func (s String) Type() string {
	return "string"
}

type Array struct {
	items []Value
}

// Hash implements Value.
func (a *Array) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: array")
}

// Index implements Indexable.
func (a *Array) Index(i int) Value {
	return a.items[i]
}

// Len implements Indexable.
func (a *Array) Len() int {
	return len(a.items)
}

// String implements Value.
func (a *Array) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, item := range a.items {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(item.String())
	}
	out.WriteString("]")
	return out.String()
}

// Truth implements Value.
func (a *Array) Truth() bool {
	return true
}

// Type implements Value.
func (a *Array) Type() string {
	return "array"
}

type MapEntry struct {
	Key   Value
	Value Value
}

type Map struct {
	entries map[uint32]MapEntry
}

// Len implements Sequence.
func (m *Map) Len() int {
	return len(m.entries)
}

// Get implements Mapping.
func (m *Map) Get(v Value) (_ Value, _ bool, err error) {
	hash, err := v.Hash()
	if err != nil {
		return nil, false, err
	}
	val, ok := m.entries[hash]
	if ok {
		return val.Value, true, nil
	}
	return Null, false, nil
}

// Hash implements Value.
func (m *Map) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: map")
}

// String implements Value.
func (m *Map) String() string {
	var entries = make([]string, 0)
	for _, item := range m.entries {
		entries = append(entries, fmt.Sprintf("%s: %s", item.Key.String(), item.Value.String()))
	}

	var out bytes.Buffer
	out.WriteString("{")
	out.WriteString(strings.Join(entries, ", "))
	out.WriteString("}")
	return out.String()
}

// Truth implements Value.
func (m *Map) Truth() bool {
	return true
}

// Type implements Value.
func (m *Map) Type() string {
	return "map"
}

type returnValue struct {
	Value Value
}

// Hash implements Value.
func (r *returnValue) Hash() (uint32, error) {
	panic("unimplemented")
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

// Hash implements Value.
func (f *Function) Hash() (uint32, error) {
	panic("unimplemented")
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

type BuiltinFunction struct {
	name string
	fn   func(args ...Value) (Value, error)
}

// Hash implements Value.
func (b *BuiltinFunction) Hash() (uint32, error) {
	panic("unimplemented")
}

// CallInternal implements Callable.
func (b *BuiltinFunction) CallInternal(args ...Value) (_ Value, err error) {
	return b.fn(args...)
}

// Name implements Callable.
func (b *BuiltinFunction) Name() string {
	return b.name
}

// String implements Value.
func (b *BuiltinFunction) String() string {
	return fmt.Sprintf("<built-in function %s>", b.Name())
}

// Truth implements Value.
func (b *BuiltinFunction) Truth() bool {
	return true
}

// Type implements Value.
func (b *BuiltinFunction) Type() string {
	return "builtin_function"
}

func NewBuiltinFunction(name string, fn func(args ...Value) (Value, error)) *BuiltinFunction {
	return &BuiltinFunction{name, fn}
}

var (
	_ Value          = NullType(0)
	_ Value          = Int(0)
	_ HasUnary       = Int(0)
	_ HasBinary      = Int(0)
	_ TotallyOrdered = Int(0)
	_ Value          = Bool(false)
	_ Comparable     = Bool(false)
	_ Value          = String("")
	_ Indexable      = String("")
	_ HasBinary      = String("")
	_ Value          = (*Array)(nil)
	_ Indexable      = (*Array)(nil)
	_ Value          = (*Map)(nil)
	_ Mapping        = (*Map)(nil)
	_ Sequence       = (*Map)(nil)
	_ Value          = (*returnValue)(nil)
	_ Value          = (*Function)(nil)
	_ Value          = (*BuiltinFunction)(nil)
	_ Callable       = (*BuiltinFunction)(nil)
)
