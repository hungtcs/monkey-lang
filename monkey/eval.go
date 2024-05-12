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
		return nil, fmt.Errorf("identifier not found: %s", node.Value)

	case *syntax.FunctionLiteral:
		return &Function{Params: node.Params, Body: node.Body, Env: env}, nil

	case *syntax.CallExpr:
		var function *Function
		if val, err := Eval(node.Function, env); err != nil {
			return nil, err
		} else if v, ok := val.(*Function); !ok {
			return nil, fmt.Errorf("not a function: %s", node.Function)
		} else {
			function = v
		}

		// 对参数求值
		args, err := evalExprs(node.Args, env)
		if err != nil {
			return nil, err
		}

		// 扩展函数 env
		fnEnv := NewEnv(env)
		for idx, param := range function.Params {
			fnEnv.Set(param.Value, args[idx])
		}
		// 执行函数体
		value, err := evalBlockStmt(function.Body, fnEnv)
		if err != nil {
			return nil, err
		}
		// 对返回值解包
		if rv, ok := value.(*returnValue); ok {
			return rv.Value, nil
		}
		return value, nil

	}
	return Null, nil
}

func evalProgram(program *syntax.Program, env *Env) (_ Value, err error) {
	var value Value
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

func Unary(op syntax.TokenType, x Value) (_ Value, err error) {
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

func Binary(op syntax.TokenType, x, y Value) (_ Value, err error) {
	if x, ok := x.(HasBinary); ok {
		v, err := x.Binary(op, y)
		if v != nil || err != nil {
			return v, err
		}
	}
	return nil, fmt.Errorf("unknown binary operator: %s %s %s", x, op, y)
}

func Compare(op syntax.TokenType, x, y Value) (_ Value, err error) {
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
