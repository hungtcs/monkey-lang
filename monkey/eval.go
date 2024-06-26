package monkey

import (
	"fmt"

	"github.com/hungtcs/monkey-lang/syntax"
)

func Eval(node syntax.Node, env *Env) (_ Value, err error) {
	switch node := node.(type) {

	case *syntax.Program:
		return evalProgram(node, env)

	case *syntax.ExprStmt:
		return Eval(node.Expr, env)

	case *syntax.IntegerLiteral:
		return Int(node.Value), nil

	case *syntax.Boolean:
		return Bool(node.Value), nil

	case *syntax.StringLiteral:
		return String(node.Value), nil

	case *syntax.ArrayLiteral:
		items, err := evalExprs(node.Items, env)
		if err != nil {
			return nil, err
		}
		return &Array{items: items}, nil

	case *syntax.MapLiteral:
		return evalMapLiteral(node, env)

	case *syntax.IndexExpr:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, err
		}
		index, err := Eval(node.Index, env)
		if err != nil {
			return nil, err
		}
		return parseIndexExpr(left, index)

	case *syntax.PrefixExpr:
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, err
		}
		return Unary(node.Op, right)

	case *syntax.InfixExpr:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, err
		}
		switch node.Op {
		case syntax.EQ, syntax.NE, syntax.GT, syntax.GE, syntax.LT, syntax.LE:
			return Compare(node.Op, left, right)
		default:
			return Binary(node.Op, left, right)
		}

	case *syntax.BlockStmt:
		return evalBlockStmt(node, env)

	case *syntax.IfExpr:
		cond, err := Eval(node.Cond, env)
		if err != nil {
			return nil, err
		}
		if cond.Truth() {
			return evalBlockStmt(node.Consequence, env)
		}
		if node.Alternative != nil {
			return evalBlockStmt(node.Alternative, env)
		}
		return Null, nil

	case *syntax.ReturnStmt:
		val, err := Eval(node.Value, env)
		if err != nil {
			return nil, err
		}
		return &returnValue{Value: val}, nil

	case *syntax.LetStmt:
		value, err := Eval(node.Value, env)
		if err != nil {
			return nil, err
		}
		env.Set(node.Name.Value, value)
		return Null, nil

	case *syntax.Identifier:
		if val, ok := env.Get(node.Value); ok {
			return val, nil
		}
		if val, ok := Universe[node.Value]; ok {
			return val, nil
		}
		return nil, fmt.Errorf("identifier not found: %s", node.Value)

	case *syntax.FunctionLiteral:
		return &Function{Params: node.Params, Body: node.Body, Env: env}, nil

	case *syntax.CallExpr:
		function, err := Eval(node.Function, env)
		if err != nil {
			return nil, err
		}

		// 对参数求值
		args, err := evalExprs(node.Args, env)
		if err != nil {
			return nil, err
		}

		// 函数调用
		return Call(function, args...)

	}
	return Null, nil
}

func evalProgram(program *syntax.Program, env *Env) (_ Value, err error) {
	var value Value = Null
	for _, stmt := range program.Stmts {
		value, err = Eval(stmt, env)
		if err != nil {
			return nil, err
		}
		// return 则提前返回，不再往后执行
		if val, ok := value.(*returnValue); ok {
			return val.Value, nil
		}
	}
	return value, nil
}

func evalBlockStmt(block *syntax.BlockStmt, env *Env) (_ Value, err error) {
	var value Value
	for _, stmt := range block.Stmts {
		value, err = Eval(stmt, env)
		if err != nil {
			return nil, err
		}

		// 如果是return语句，则跳出所有代码块
		if _, ok := value.(*returnValue); ok {
			return value, nil
		}
	}
	return value, nil
}

func evalExprs(exprs []syntax.Expr, env *Env) (_ []Value, err error) {
	var values = make([]Value, len(exprs))
	for i, expr := range exprs {
		values[i], err = Eval(expr, env)
		if err != nil {
			return nil, err
		}
	}
	return values, nil
}

func evalMapLiteral(node *syntax.MapLiteral, env *Env) (_ Value, err error) {
	entries := make(map[uint32]MapEntry)
	for keyNode, valNode := range node.Pairs {
		key, err := Eval(keyNode, env)
		if err != nil {
			return nil, err
		}
		hash, err := key.Hash()
		if err != nil {
			return nil, err
		}
		val, err := Eval(valNode, env)
		if err != nil {
			return nil, err
		}
		entries[hash] = MapEntry{Key: key, Value: val}
	}
	return &Map{entries: entries}, nil
}

func parseIndexExpr(left, index Value) (_ Value, err error) {
	switch left := left.(type) {
	case Mapping:
		val, found, err := left.Get(index)
		if err != nil {
			return nil, err
		} else if found {
			return val, nil
		} else {
			return Null, nil
		}
	case Indexable:
		if iv, ok := index.(Int); !ok {
			return nil, fmt.Errorf("invalid index type: %s", index.Type())
		} else {
			// 支持负数索引
			if iv < 0 {
				iv = Int(left.Len()) + (iv)
			}
			return left.Index(int(iv)), nil
		}
	}
	return nil, fmt.Errorf("index operator not supported: %s", left.Type())
}

func Unary(op syntax.Token, x Value) (_ Value, err error) {
	if op == syntax.BANG {
		return Bool(!x.Truth()), nil
	}

	if x, ok := x.(HasUnary); ok {
		v, err := x.Unary(op)
		if v != nil || err != nil {
			return v, err
		}
	}

	return nil, fmt.Errorf("unknown unary operator: %s", op)
}

func Binary(op syntax.Token, x, y Value) (_ Value, err error) {
	if x, ok := x.(HasBinary); ok {
		v, err := x.Binary(op, y, Left)
		if v != nil || err != nil {
			return v, err
		}
	}

	if y, ok := y.(HasBinary); ok {
		v, err := y.Binary(op, x, Right)
		if v != nil || err != nil {
			return v, err
		}
	}
	return nil, fmt.Errorf("unknown binary operator: %s %s %s", x, op, y)
}

func Compare(op syntax.Token, x, y Value) (_ Value, err error) {
	if isSameType(x, y) {
		if x, ok := x.(Comparable); ok {
			return x.Compare(op, y)
		}
		if xcomp, ok := x.(TotallyOrdered); ok {
			t, err := xcomp.Cmp(y)
			if err != nil {
				return False, err
			}
			return Bool(threeway(op, t)), nil
		}
	}
	return nil, fmt.Errorf("invalid cmp operator: %s %s %s", x, syntax.EQ, y)
}

func Call(value Value, args ...Value) (_ Value, err error) {
	switch value := value.(type) {
	case *Function:
		// 扩展函数 env
		fnEnv := NewEnv(value.Env)
		for idx, param := range value.Params {
			fnEnv.Set(param.Value, args[idx])
		}
		// 执行函数体
		result, err := evalBlockStmt(value.Body, fnEnv)
		if err != nil {
			return nil, err
		}
		// 对返回值解包
		if rv, ok := result.(*returnValue); ok {
			return rv.Value, nil
		}
		return value, nil

	case *BuiltinFunction:
		return value.CallInternal(args...)

	}
	return nil, fmt.Errorf("invalid call of non-function (%s)", value.Type())
}
