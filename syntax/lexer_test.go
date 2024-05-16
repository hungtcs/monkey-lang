package syntax

import (
	"fmt"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
`

	tests := []struct {
		expectedType    Token
		expectedLiteral string
		pos             string
	}{
		{LET, "let", "1:1"},
		{IDENT, "five", "1:5"},
		{ASSIGN, "=", "1:10"},
		{INT, "5", "1:11"},
		{SEMICOLON, ";", "1:13"},
		{LET, "let", "2:1"},
		{IDENT, "ten", "2:5"},
		{ASSIGN, "=", "2:9"},
		{INT, "10", "2:11"},
		{SEMICOLON, ";", "2:14"},
		{LET, "let", "3:1"},
		{IDENT, "add", "4:5"},
		{ASSIGN, "=", "4:9"},
		{FUNCTION, "fn", "4:11"},
		{LPAREN, "(", "4:14"},
		{IDENT, "x", "4:15"},
		{COMMA, ",", "4:16"},
		{IDENT, "y", "4:17"},
		{RPAREN, ")", "4:19"},
		{LBRACE, "{", "4:20"},
		{IDENT, "x", "5:1"},
		{PLUS, "+", "5:4"},
		{IDENT, "y", "5:6"},
		{SEMICOLON, ";", "5:8"},
		{RBRACE, "}", "6:1"},
		{SEMICOLON, ";", "6:3"},
		{LET, "let", "7:1"},
		{IDENT, "result", "8:5"},
		{ASSIGN, "=", "8:12"},
		{IDENT, "add", "8:14"},
		{LPAREN, "(", "8:18"},
		{IDENT, "five", "8:19"},
		{COMMA, ",", "8:23"},
		{IDENT, "ten", "8:24"},
		{RPAREN, ")", "8:28"},
		{SEMICOLON, ";", "8:29"},
		{BANG, "!", "9:1"},
		{MINUS, "-", "9:3"},
		{SLASH, "/", "9:4"},
		{STAR, "*", "9:5"},
		{INT, "5", "9:6"},
		{SEMICOLON, ";", "9:7"},
		{INT, "5", "10:1"},
		{LT, "<", "10:3"},
		{INT, "10", "10:5"},
		{GT, ">", "10:8"},
		{INT, "5", "10:10"},
		{SEMICOLON, ";", "10:12"},
		{IF, "if", "11:1"},
		{LPAREN, "(", "12:4"},
		{INT, "5", "12:6"},
		{LT, "<", "12:7"},
		{INT, "10", "12:9"},
		{RPAREN, ")", "12:12"},
		{LBRACE, "{", "12:13"},
		{RETURN, "return", "13:1"},
		{TRUE, "true", "13:9"},
		{SEMICOLON, ";", "13:14"},
		{RBRACE, "}", "14:1"},
		{ELSE, "else", "14:3"},
		{LBRACE, "{", "14:8"},
		{RETURN, "return", "15:1"},
		{FALSE, "false", "15:9"},
		{SEMICOLON, ";", "15:15"},
		{RBRACE, "}", "16:1"},
		{INT, "10", "17:1"},
		{EQ, "==", "18:4"},
		{INT, "10", "18:7"},
		{SEMICOLON, ";", "18:10"},
		{INT, "10", "19:1"},
		{NE, "!=", "19:4"},
		{INT, "9", "19:7"},
		{SEMICOLON, ";", "19:9"},
		{STRING, "foobar", "20:1"},
		{STRING, "foo bar", "21:1"},
		{LBRACKET, "[", "22:1"},
		{INT, "1", "22:3"},
		{COMMA, ",", "22:4"},
		{INT, "2", "22:5"},
		{RBRACKET, "]", "22:7"},
		{SEMICOLON, ";", "22:8"},
		{LBRACE, "{", "23:1"},
		{STRING, "foo", "23:3"},
		{COLON, ":", "23:8"},
		{STRING, "bar", "23:9"},
		{RBRACE, "}", "23:15"},
		{EOF, "", "24:1"},
	}

	lexer := NewLexer(input)
	for i, tt := range tests {
		tok := lexer.NextToken()
		pos := tok.pos

		fmt.Printf("%v \t\t %v\n", pos, tok)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
		// if fmt.Sprintf("%d:%d", pos.Line, pos.Col) != tt.pos {
		// 	t.Fatalf("tests[%d] - literal wrong. expected=%s, got=%s",
		// 		i, tt.pos, fmt.Sprintf("%d:%d", pos.Line, pos.Col))
		// }
	}
}
