package phase1

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;
let add = fn(x: int, y: int) -> int {
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
while (true) {
    break;
}
for i in 0..10 {
    continue;
}
3.14
let pi: float = 3.14159;
// これはコメントです
let flag: bool = true;`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LET, "let"},
		{IDENT, "five"},
		{ASSIGN, "="},
		{INT, "5"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "ten"},
		{ASSIGN, "="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "add"},
		{ASSIGN, "="},
		{FN, "fn"},
		{LPAREN, "("},
		{IDENT, "x"},
		{COLON, ":"},
		{INT_TYPE, "int"},
		{COMMA, ","},
		{IDENT, "y"},
		{COLON, ":"},
		{INT_TYPE, "int"},
		{RPAREN, ")"},
		{ARROW, "->"},
		{INT_TYPE, "int"},
		{LBRACE, "{"},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "five"},
		{COMMA, ","},
		{IDENT, "ten"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{NOT, "!"},
		{MINUS, "-"},
		{DIVIDE, "/"},
		{MULTIPLY, "*"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{GT, ">"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{IF, "if"},
		{LPAREN, "("},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{ELSE, "else"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{FALSE, "false"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{INT, "10"},
		{EQ, "=="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{INT, "10"},
		{NOT_EQ, "!="},
		{INT, "9"},
		{SEMICOLON, ";"},
		{STRING, "foobar"},
		{STRING, "foo bar"},
		{LBRACKET, "["},
		{INT, "1"},
		{COMMA, ","},
		{INT, "2"},
		{RBRACKET, "]"},
		{SEMICOLON, ";"},
		{LBRACE, "{"},
		{STRING, "foo"},
		{COLON, ":"},
		{STRING, "bar"},
		{RBRACE, "}"},
		{WHILE, "while"},
		{LPAREN, "("},
		{TRUE, "true"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{BREAK, "break"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{FOR, "for"},
		{IDENT, "i"},
		{IDENT, "in"},
		{INT, "0"},
		{IDENT, ".."},
		{INT, "10"},
		{LBRACE, "{"},
		{CONTINUE, "continue"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{FLOAT, "3.14"},
		{LET, "let"},
		{IDENT, "pi"},
		{COLON, ":"},
		{FLOAT_TYPE, "float"},
		{ASSIGN, "="},
		{FLOAT, "3.14159"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "flag"},
		{COLON, ":"},
		{BOOL_TYPE, "bool"},
		{ASSIGN, "="},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIdentifiers(t *testing.T) {
	tests := []struct {
		input    string
		expected []struct {
			tokenType TokenType
			literal   string
		}
	}{
		{
			"identifier",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{IDENT, "identifier"},
			},
		},
		{
			"_private",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{IDENT, "_private"},
			},
		},
		{
			"camelCase",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{IDENT, "camelCase"},
			},
		},
		{
			"snake_case",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{IDENT, "snake_case"},
			},
		},
		{
			"variable123",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{IDENT, "variable123"},
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)
		for i, expected := range tt.expected {
			tok := l.NextToken()
			if tok.Type != expected.tokenType {
				t.Errorf("test[%d] wrong token type. expected=%q, got=%q", i, expected.tokenType, tok.Type)
			}
			if tok.Literal != expected.literal {
				t.Errorf("test[%d] wrong literal. expected=%q, got=%q", i, expected.literal, tok.Literal)
			}
		}
	}
}

func TestNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected []struct {
			tokenType TokenType
			literal   string
		}
	}{
		{
			"123",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{INT, "123"},
			},
		},
		{
			"0",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{INT, "0"},
			},
		},
		{
			"3.14",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{FLOAT, "3.14"},
			},
		},
		{
			"0.5",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{FLOAT, "0.5"},
			},
		},
		{
			"123.456",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{FLOAT, "123.456"},
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)
		for i, expected := range tt.expected {
			tok := l.NextToken()
			if tok.Type != expected.tokenType {
				t.Errorf("test[%d] wrong token type. expected=%q, got=%q", i, expected.tokenType, tok.Type)
			}
			if tok.Literal != expected.literal {
				t.Errorf("test[%d] wrong literal. expected=%q, got=%q", i, expected.literal, tok.Literal)
			}
		}
	}
}

func TestStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected []struct {
			tokenType TokenType
			literal   string
		}
	}{
		{
			`"hello"`,
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{STRING, "hello"},
			},
		},
		{
			`"hello world"`,
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{STRING, "hello world"},
			},
		},
		{
			`"hello\nworld"`,
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{STRING, "hello\nworld"},
			},
		},
		{
			`"hello \"world\""`,
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{STRING, `hello \"world\"`},
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)
		for i, expected := range tt.expected {
			tok := l.NextToken()
			if tok.Type != expected.tokenType {
				t.Errorf("test[%d] wrong token type. expected=%q, got=%q", i, expected.tokenType, tok.Type)
			}
			if tok.Literal != expected.literal {
				t.Errorf("test[%d] wrong literal. expected=%q, got=%q", i, expected.literal, tok.Literal)
			}
		}
	}
}

func TestComments(t *testing.T) {
	input := `let x = 5; // これはコメントです
// 完全にコメントの行
let y = 10;`

	expected := []struct {
		tokenType TokenType
		literal   string
	}{
		{LET, "let"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{INT, "5"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "y"},
		{ASSIGN, "="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	l := New(input)
	for i, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.tokenType {
			t.Errorf("test[%d] wrong token type. expected=%q, got=%q", i, tt.tokenType, tok.Type)
		}
		if tok.Literal != tt.literal {
			t.Errorf("test[%d] wrong literal. expected=%q, got=%q", i, tt.literal, tok.Literal)
		}
	}
}

func TestLineAndColumnNumbers(t *testing.T) {
	input := `let x = 5;
let y = 10;`

	expected := []struct {
		tokenType TokenType
		literal   string
		line      int
		column    int
	}{
		{LET, "let", 1, 1},
		{IDENT, "x", 1, 5},
		{ASSIGN, "=", 1, 7},
		{INT, "5", 1, 9},
		{SEMICOLON, ";", 1, 10},
		{LET, "let", 2, 1},
		{IDENT, "y", 2, 5},
		{ASSIGN, "=", 2, 7},
		{INT, "10", 2, 9},
		{SEMICOLON, ";", 2, 11},
		{EOF, "", 2, 12},
	}

	l := New(input)
	for i, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.tokenType {
			t.Errorf("test[%d] wrong token type. expected=%q, got=%q", i, tt.tokenType, tok.Type)
		}
		if tok.Literal != tt.literal {
			t.Errorf("test[%d] wrong literal. expected=%q, got=%q", i, tt.literal, tok.Literal)
		}
		if tok.Line != tt.line {
			t.Errorf("test[%d] wrong line. expected=%d, got=%d", i, tt.line, tok.Line)
		}
		if tok.Column != tt.column {
			t.Errorf("test[%d] wrong column. expected=%d, got=%d", i, tt.column, tok.Column)
		}
	}
}

func TestIllegalCharacters(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    TokenType
		expectedLiteral string
	}{
		{"@", ILLEGAL, "@"},
		{"#", ILLEGAL, "#"},
		{"$", ILLEGAL, "$"},
		{"&", ILLEGAL, "&"}, // 単体の&は不正（&&は論理演算子）
		{"|", ILLEGAL, "|"}, // 単体の|は不正（||は論理演算子）
	}

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] wrong token type. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] wrong literal. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestKeywordRecognition(t *testing.T) {
	tests := []struct {
		input        string
		expectedType TokenType
	}{
		{"let", LET},
		{"fn", FN},
		{"if", IF},
		{"else", ELSE},
		{"return", RETURN},
		{"true", TRUE},
		{"false", FALSE},
		{"while", WHILE},
		{"for", FOR},
		{"break", BREAK},
		{"continue", CONTINUE},
		{"int", INT_TYPE},
		{"float", FLOAT_TYPE},
		{"string", STRING_TYPE},
		{"bool", BOOL_TYPE},
		{"notakeyword", IDENT}, // キーワードではない
	}

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] wrong token type. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
	}
}
