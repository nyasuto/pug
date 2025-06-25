package phase1

import (
	"testing"
)

// TestParserErrors はパーサーのエラーハンドリングをテストする
func TestParserErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "missing identifier in let statement",
			input:         "let = 5;",
			expectedError: "expected next token to be IDENT, got = instead",
		},
		{
			name:          "missing assign in let statement",
			input:         "let x 5;",
			expectedError: "expected next token to be =, got INT instead",
		},
		{
			name:          "invalid integer literal",
			input:         "let x = 999999999999999999999999999999999999999999999999999999999999999999;",
			expectedError: "could not parse \"999999999999999999999999999999999999999999999999999999999999999999\" as integer",
		},
		{
			name:          "missing closing parenthesis in grouped expression",
			input:         "let x = (5 + 3;",
			expectedError: "expected next token to be ), got ; instead",
		},
		{
			name:          "missing opening parenthesis in if expression",
			input:         "if true { 5 }",
			expectedError: "expected next token to be (, got TRUE instead",
		},
		{
			name:          "missing closing parenthesis in if expression",
			input:         "if (true { 5 }",
			expectedError: "expected next token to be ), got { instead",
		},
		{
			name:          "missing opening brace in if expression",
			input:         "if (true) 5 }",
			expectedError: "expected next token to be {, got INT instead",
		},
		{
			name:          "missing opening brace in else",
			input:         "if (true) { 5 } else 3 }",
			expectedError: "expected next token to be {, got INT instead",
		},
		{
			name:          "missing opening parenthesis in function",
			input:         "fn x, y { x + y }",
			expectedError: "expected next token to be (, got IDENT instead",
		},
		{
			name:          "missing closing parenthesis in function",
			input:         "fn (x, y { x + y }",
			expectedError: "expected next token to be ), got { instead",
		},
		{
			name:          "missing opening brace in function",
			input:         "fn (x, y) x + y }",
			expectedError: "expected next token to be {, got IDENT instead",
		},
		{
			name:          "missing closing parenthesis in call expression",
			input:         "add(1, 2",
			expectedError: "expected next token to be ), got EOF instead",
		},
		{
			name:          "missing opening parenthesis in while statement",
			input:         "while true { break; }",
			expectedError: "expected next token to be (, got TRUE instead",
		},
		{
			name:          "missing closing parenthesis in while statement",
			input:         "while (true { break; }",
			expectedError: "expected next token to be ), got { instead",
		},
		{
			name:          "missing opening brace in while statement",
			input:         "while (true) break; }",
			expectedError: "expected next token to be {, got BREAK instead",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			p := NewParser(l)
			p.ParseProgram()

			errors := p.Errors()
			if len(errors) == 0 {
				t.Errorf("expected parser errors, but got none")
				return
			}

			found := false
			for _, err := range errors {
				if err == tt.expectedError {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected error %q not found in errors: %v", tt.expectedError, errors)
			}
		})
	}
}

// TestNoPrefixParseFnError は前置構文解析関数が見つからない場合のエラーをテストする
func TestNoPrefixParseFnError(t *testing.T) {
	input := "5 + +;"

	l := New(input)
	p := NewParser(l)
	p.ParseProgram()

	errors := p.Errors()
	if len(errors) == 0 {
		t.Errorf("expected parser errors, but got none")
		return
	}

	expectedError := "no prefix parse function for ; found"
	found := false
	for _, err := range errors {
		if err == expectedError {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected error %q not found in errors: %v", expectedError, errors)
	}
}

// TestParseEmptyInput は空の入力に対するパーサーの動作をテストする
func TestParseEmptyInput(t *testing.T) {
	input := ""

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) != 0 {
		t.Errorf("unexpected parser errors: %v", errors)
	}

	if program == nil {
		t.Errorf("ParseProgram() returned nil")
		return
	}

	if len(program.Statements) != 0 {
		t.Errorf("expected empty program, got %d statements", len(program.Statements))
	}
}

// TestParseInvalidTokens は無効なトークンに対するパーサーの動作をテストする
func TestParseInvalidTokens(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unknown character",
			input: "@",
		},
		{
			name:  "incomplete string",
			input: "\"hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			p := NewParser(l)
			program := p.ParseProgram()

			// パーサーがクラッシュしないことを確認
			if program == nil {
				t.Errorf("ParseProgram() returned nil")
			}
		})
	}
}

// TestExpressionStatementWithNil は式がnilの場合のExpressionStatementをテストする
func TestExpressionStatementWithNil(t *testing.T) {
	// 無効な式を作成してnilのExpressionを持つExpressionStatementをテストする
	stmt := &ExpressionStatement{
		Token:      Token{Type: INT, Literal: ""},
		Expression: nil,
	}

	result := stmt.String()
	expected := ""

	if result != expected {
		t.Errorf("ExpressionStatement.String() with nil expression = %q, want %q", result, expected)
	}
}

// TestReturnStatementWithNilValue はnilの戻り値を持つReturnStatementをテストする
func TestReturnStatementWithNilValue(t *testing.T) {
	stmt := &ReturnStatement{
		Token:       Token{Type: RETURN, Literal: "return"},
		ReturnValue: nil,
	}

	result := stmt.String()
	expected := "return ;"

	if result != expected {
		t.Errorf("ReturnStatement.String() with nil value = %q, want %q", result, expected)
	}
}

// TestLetStatementWithNilValue はnilの値を持つLetStatementをテストする
func TestLetStatementWithNilValue(t *testing.T) {
	stmt := &LetStatement{
		Token: Token{Type: LET, Literal: "let"},
		Name: &Identifier{
			Token: Token{Type: IDENT, Literal: "x"},
			Value: "x",
		},
		Value: nil,
	}

	result := stmt.String()
	expected := "let x = ;"

	if result != expected {
		t.Errorf("LetStatement.String() with nil value = %q, want %q", result, expected)
	}
}
