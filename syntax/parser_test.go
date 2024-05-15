package syntax

import (
	"fmt"
	"testing"
)

func TestLetStmts(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		p := NewParser(tt.input)
		program, err := p.Parse()
		checkParserErrors(t, err)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain 1 Stmts. got=%d",
				len(program.Stmts))
		}

		stmt := program.Stmts[0]
		if !testLetStmt(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*LetStmt).Value
		if !testLiteralExpr(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStmts(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		p := NewParser(tt.input)
		program, err := p.Parse()
		checkParserErrors(t, err)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain 1 Stmts. got=%d",
				len(program.Stmts))
		}

		stmt := program.Stmts[0]
		returnStmt, ok := stmt.(*ReturnStmt)
		if !ok {
			t.Fatalf("stmt not *ReturnStatement. got=%T", stmt)
		}
		if returnStmt.Literal() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.Literal())
		}
		if testLiteralExpr(t, returnStmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpr(t *testing.T) {
	input := "foobar;"

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	if len(program.Stmts) != 1 {
		t.Fatalf("program has not enough Stmts. got=%d",
			len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ExprStmt. got=%T",
			program.Stmts[0])
	}

	ident, ok := stmt.Expr.(*Identifier)
	if !ok {
		t.Fatalf("exp not *Identifier. got=%T", stmt.Expr)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.Literal() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.Literal())
	}
}

func TestIntegerLiteralExpr(t *testing.T) {
	input := "5;"

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	if len(program.Stmts) != 1 {
		t.Fatalf("program has not enough Stmts. got=%d",
			len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ExprStmt. got=%T",
			program.Stmts[0])
	}

	literal, ok := stmt.Expr.(*IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *IntegerLiteral. got=%T", stmt.Expr)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.Literal() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.Literal())
	}
}

func TestParsingPrefixExprs(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		// l := lexer.New(tt.input)
		p := NewParser(tt.input)
		program, err := p.Parse()
		checkParserErrors(t, err)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain %d Stmts. got=%d\n",
				1, len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ExprStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not ExprStmt. got=%T",
				program.Stmts[0])
		}

		exp, ok := stmt.Expr.(*PrefixExpr)
		if !ok {
			t.Fatalf("stmt is not PrefixExpr. got=%T", stmt.Expr)
		}
		if exp.Op.String() != tt.operator {
			t.Fatalf("exp.Op is not '%s'. got=%s",
				tt.operator, exp.Op)
		}
		if !testLiteralExpr(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExprs(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		// l := lexer.New(tt.input)
		p := NewParser(tt.input)
		program, err := p.Parse()
		checkParserErrors(t, err)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain %d Stmts. got=%d\n",
				1, len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ExprStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not ExprStmt. got=%T",
				program.Stmts[0])
		}

		if !testInfixExpr(t, stmt.Expr, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		// {
		// 	"!-a",
		// 	"(!(-a))",
		// },
		// {
		// 	"a + b + c",
		// 	"((a + b) + c)",
		// },
		// {
		// 	"a + b - c",
		// 	"((a + b) - c)",
		// },
		// {
		// 	"a * b * c",
		// 	"((a * b) * c)",
		// },
		// {
		// 	"a * b / c",
		// 	"((a * b) / c)",
		// },
		// {
		// 	"a + b / c",
		// 	"(a + (b / c))",
		// },
		// {
		// 	"a + b * c + d / e - f",
		// 	"(((a + (b * c)) + (d / e)) - f)",
		// },
		// {
		// 	"3 + 4; -5 * 5",
		// 	"(3 + 4)((-5) * 5)",
		// },
		// {
		// 	"5 > 4 == 3 < 4",
		// 	"((5 > 4) == (3 < 4))",
		// },
		// {
		// 	"5 < 4 != 3 > 4",
		// 	"((5 < 4) != (3 > 4))",
		// },
		// {
		// 	"3 + 4 * 5 == 3 * 1 + 4 * 5",
		// 	"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		// },
		// {
		// 	"true",
		// 	"true",
		// },
		// {
		// 	"false",
		// 	"false",
		// },
		// {
		// 	"3 > 5 == false",
		// 	"((3 > 5) == false)",
		// },
		// {
		// 	"3 < 5 == true",
		// 	"((3 < 5) == true)",
		// },
		// {
		// 	"1 + (2 + 3) + 4",
		// 	"((1 + (2 + 3)) + 4)",
		// },
		// {
		// 	"(5 + 5) * 2",
		// 	"((5 + 5) * 2)",
		// },
		// {
		// 	"2 / (5 + 5)",
		// 	"(2 / (5 + 5))",
		// },
		// {
		// 	"(5 + 5) * 2 * (5 + 5)",
		// 	"(((5 + 5) * 2) * (5 + 5))",
		// },
		// {
		// 	"-(5 + 5)",
		// 	"(-(5 + 5))",
		// },
		// {
		// 	"!(true == true)",
		// 	"(!(true == true))",
		// },
		// {
		// 	"a + add(b * c) + d",
		// 	"((a + add((b * c))) + d)",
		// },
		// {
		// 	"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
		// 	"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		// },
		// {
		// 	"add(a + b + c * d / f + g)",
		// 	"add((((a + b) + ((c * d) / f)) + g))",
		// },
		// {
		// 	"a * [1, 2, 3, 4][b * c] * d",
		// 	"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		// },
		// {
		// 	"add(a * b[2], b[1], 2 * [1, 2][1])",
		// 	"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		// },
	}

	for _, tt := range tests {
		// l := lexer.New(tt.input)
		p := NewParser(tt.input)
		program, err := p.Parse()
		checkParserErrors(t, err)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpr(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		// l := lexer.New(tt.input)
		p := NewParser(tt.input)
		program, err := p.Parse()
		checkParserErrors(t, err)

		if len(program.Stmts) != 1 {
			t.Fatalf("program has not enough Stmts. got=%d",
				len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ExprStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not ExprStmt. got=%T",
				program.Stmts[0])
		}

		boolean, ok := stmt.Expr.(*Boolean)
		if !ok {
			t.Fatalf("exp not *Boolean. got=%T", stmt.Expr)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestIfExpr(t *testing.T) {
	input := `if (x < y) { x }`

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	if len(program.Stmts) != 1 {
		t.Fatalf("program.Stmts does not contain %d Stmts. got=%d\n",
			1, len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ExprStmt. got=%T",
			program.Stmts[0])
	}

	exp, ok := stmt.Expr.(*IfExpr)
	if !ok {
		t.Fatalf("stmt.Expr is not IfExpr. got=%T",
			stmt.Expr)
	}

	if !testInfixExpr(t, exp.Cond, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Stmts) != 1 {
		t.Errorf("consequence is not 1 Stmts. got=%d\n",
			len(exp.Consequence.Stmts))
	}

	consequence, ok := exp.Consequence.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("Stmts[0] is not ExprStmt. got=%T",
			exp.Consequence.Stmts[0])
	}

	if !testIdentifier(t, consequence.Expr, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Stmts was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpr(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	if len(program.Stmts) != 1 {
		t.Fatalf("program.Stmts does not contain %d Stmts. got=%d\n",
			1, len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ExprStmt. got=%T",
			program.Stmts[0])
	}

	exp, ok := stmt.Expr.(*IfExpr)
	if !ok {
		t.Fatalf("stmt.Expr is not IfExpr. got=%T", stmt.Expr)
	}

	if !testInfixExpr(t, exp.Cond, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Stmts) != 1 {
		t.Errorf("consequence is not 1 Stmts. got=%d\n",
			len(exp.Consequence.Stmts))
	}

	consequence, ok := exp.Consequence.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("Stmts[0] is not ExprStmt. got=%T",
			exp.Consequence.Stmts[0])
	}

	if !testIdentifier(t, consequence.Expr, "x") {
		return
	}

	if len(exp.Alternative.Stmts) != 1 {
		t.Errorf("exp.Alternative.Stmts does not contain 1 Stmts. got=%d\n",
			len(exp.Alternative.Stmts))
	}

	alternative, ok := exp.Alternative.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("Stmts[0] is not ExprStmt. got=%T",
			exp.Alternative.Stmts[0])
	}

	if !testIdentifier(t, alternative.Expr, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	if len(program.Stmts) != 1 {
		t.Fatalf("program.Stmts does not contain %d Stmts. got=%d\n",
			1, len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ExprStmt. got=%T",
			program.Stmts[0])
	}

	function, ok := stmt.Expr.(*FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expr is not FunctionLiteral. got=%T",
			stmt.Expr)
	}

	if len(function.Params) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Params))
	}

	testLiteralExpr(t, function.Params[0], "x")
	testLiteralExpr(t, function.Params[1], "y")

	if len(function.Body.Stmts) != 1 {
		t.Fatalf("function.Body.Stmts has not 1 Stmts. got=%d\n",
			len(function.Body.Stmts))
	}

	bodyStmt, ok := function.Body.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("function body stmt is not ExprStmt. got=%T",
			function.Body.Stmts[0])
	}

	testInfixExpr(t, bodyStmt.Expr, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		// l := lexer.New(tt.input)
		p := NewParser(tt.input)
		program, err := p.Parse()
		checkParserErrors(t, err)

		stmt := program.Stmts[0].(*ExprStmt)
		function := stmt.Expr.(*FunctionLiteral)

		if len(function.Params) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Params))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpr(t, function.Params[i], ident)
		}
	}
}

func TestCallExprParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	if len(program.Stmts) != 1 {
		t.Fatalf("program.Stmts does not contain %d Stmts. got=%d\n",
			1, len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("stmt is not ExprStmt. got=%T",
			program.Stmts[0])
	}

	exp, ok := stmt.Expr.(*CallExpr)
	if !ok {
		t.Fatalf("stmt.Expr is not CallExpr. got=%T",
			stmt.Expr)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Args) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Args))
	}

	testLiteralExpr(t, exp.Args[0], 1)
	testInfixExpr(t, exp.Args[1], 2, "*", 3)
	testInfixExpr(t, exp.Args[2], 4, "+", 5)
}

func TestCallExprParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		// l := lexer.New(tt.input)
		p := NewParser(tt.input)
		program, err := p.Parse()
		checkParserErrors(t, err)

		stmt := program.Stmts[0].(*ExprStmt)
		exp, ok := stmt.Expr.(*CallExpr)
		if !ok {
			t.Fatalf("stmt.Expr is not CallExpr. got=%T",
				stmt.Expr)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Args) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Args))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Args[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Args[i].String())
			}
		}
	}
}

func TestStringLiteralExpr(t *testing.T) {
	input := `"hello world";`
	fmt.Printf("input: %v\n", input)

	// l := lexer.New(input)
	p := NewParser(input)
	fmt.Printf("p: %v\n", p)
	program, err := p.Parse()
	checkParserErrors(t, err)
	fmt.Printf("program: %v\n", program)

	stmt := program.Stmts[0].(*ExprStmt)
	literal, ok := stmt.Expr.(*StringLiteral)
	if !ok {
		t.Fatalf("exp not *StringLiteral. got=%T", stmt.Expr)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingEmptyArrayLiterals(t *testing.T) {
	input := "[]"

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	stmt, ok := program.Stmts[0].(*ExprStmt)
	array, ok := stmt.Expr.(*ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ArrayLiteral. got=%T", stmt.Expr)
	}

	if len(array.Items) != 0 {
		t.Errorf("len(array.Items) not 0. got=%d", len(array.Items))
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	stmt, ok := program.Stmts[0].(*ExprStmt)
	array, ok := stmt.Expr.(*ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ArrayLiteral. got=%T", stmt.Expr)
	}

	if len(array.Items) != 3 {
		t.Fatalf("len(array.Items) not 3. got=%d", len(array.Items))
	}

	testIntegerLiteral(t, array.Items[0], 1)
	testInfixExpr(t, array.Items[1], 2, "*", 2)
	testInfixExpr(t, array.Items[2], 3, "+", 3)
}

func TestParsingIndexExprs(t *testing.T) {
	input := "myArray[1 + 1]"

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	stmt, ok := program.Stmts[0].(*ExprStmt)
	indexExp, ok := stmt.Expr.(*IndexExpr)
	if !ok {
		t.Fatalf("exp not *IndexExpr. got=%T", stmt.Expr)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpr(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingEmptyMapLiteral(t *testing.T) {
	input := "{}"

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	stmt := program.Stmts[0].(*ExprStmt)
	hash, ok := stmt.Expr.(*MapLiteral)
	if !ok {
		t.Fatalf("exp is not MapLiteral. got=%T", stmt.Expr)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingMapLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	stmt := program.Stmts[0].(*ExprStmt)
	hash, ok := stmt.Expr.(*MapLiteral)
	if !ok {
		t.Fatalf("exp is not MapLiteral. got=%T", stmt.Expr)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*StringLiteral)
		if !ok {
			t.Errorf("key is not StringLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingMapLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	stmt := program.Stmts[0].(*ExprStmt)
	hash, ok := stmt.Expr.(*MapLiteral)
	if !ok {
		t.Fatalf("exp is not MapLiteral. got=%T", stmt.Expr)
	}

	expected := map[string]int64{
		"true":  1,
		"false": 2,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		boolean, ok := key.(*Boolean)
		if !ok {
			t.Errorf("key is not BooleanLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[boolean.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingMapLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	stmt := program.Stmts[0].(*ExprStmt)
	hash, ok := stmt.Expr.(*MapLiteral)
	if !ok {
		t.Fatalf("exp is not MapLiteral. got=%T", stmt.Expr)
	}

	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		integer, ok := key.(*IntegerLiteral)
		if !ok {
			t.Errorf("key is not IntegerLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[integer.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingMapLiteralsWithExprs(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	// l := lexer.New(input)
	p := NewParser(input)
	program, err := p.Parse()
	checkParserErrors(t, err)

	stmt := program.Stmts[0].(*ExprStmt)
	hash, ok := stmt.Expr.(*MapLiteral)
	if !ok {
		t.Fatalf("exp is not MapLiteral. got=%T", stmt.Expr)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(Expr){
		"one": func(e Expr) {
			testInfixExpr(t, e, 0, "+", 1)
		},
		"two": func(e Expr) {
			testInfixExpr(t, e, 10, "-", 8)
		},
		"three": func(e Expr) {
			testInfixExpr(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*StringLiteral)
		if !ok {
			t.Errorf("key is not StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
	}
}

func testLetStmt(t *testing.T, s Stmt, name string) bool {
	if s.Literal() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.Literal())
		return false
	}

	letStmt, ok := s.(*LetStmt)
	if !ok {
		t.Errorf("s not *LetStmt. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.Literal() != name {
		t.Errorf("letStmt.Name.Literal() not '%s'. got=%s",
			name, letStmt.Name.Literal())
		return false
	}

	return true
}

func testInfixExpr(t *testing.T, exp Expr, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*InfixExpr)
	if !ok {
		t.Errorf("exp is not InfixExpr. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpr(t, opExp.Left, left) {
		return false
	}

	if opExp.Op.String() != operator {
		t.Errorf("exp.Op is not '%s'. got=%q", operator, opExp.Op)
		return false
	}

	if !testLiteralExpr(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpr(
	t *testing.T,
	exp Expr,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il Expr, value int64) bool {
	integ, ok := il.(*IntegerLiteral)
	if !ok {
		t.Errorf("il not *IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.Literal() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.Literal())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp Expr, value string) bool {
	ident, ok := exp.(*Identifier)
	if !ok {
		t.Errorf("exp not *Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.Literal() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.Literal())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp Expr, value bool) bool {
	bo, ok := exp.(*Boolean)
	if !ok {
		t.Errorf("exp not *Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.Literal() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.Literal())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, err error) {
	if err == nil {
		return
	}

	t.Errorf("parser error: %s", err.Error())
	t.FailNow()
}
