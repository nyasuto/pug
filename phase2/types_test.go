package phase2

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nyasuto/pug/phase1"
)

// parseProgramForTypes はテスト用のヘルパー関数
func parseProgramForTypes(t *testing.T, input string) *phase1.Program {
	lexer := phase1.New(input)
	parser := phase1.NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		t.Fatalf("parser errors: %v", parser.Errors())
	}

	return program
}

// TestTypeChecker_BasicTypes は基本型の型検査をテストする
func TestTypeChecker_BasicTypes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"42", "int"},
		{"3.14", "float"},
		{"\"hello\"", "string"},
		{"true", "bool"},
		{"false", "bool"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgramForTypes(t, tt.input)
			tc := NewTypeChecker()

			resultType, errors := tc.CheckProgram(program)
			if len(errors) > 0 {
				t.Fatalf("unexpected errors: %v", errors)
			}

			if resultType.String() != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, resultType.String())
			}
		})
	}
}

// TestTypeChecker_Variables は変数の型検査をテストする
func TestTypeChecker_Variables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "let文で変数定義",
			input:    "let x = 42; x",
			expected: "int",
			hasError: false,
		},
		{
			name:     "未定義変数の参照",
			input:    "undefined_var",
			expected: "?error",
			hasError: true,
		},
		{
			name:     "複数の変数定義と使用",
			input:    "let a = 10; let b = 20; a + b",
			expected: "int",
			hasError: false,
		},
		{
			name:     "異なる型の変数",
			input:    "let x = 42; let y = \"hello\"; y",
			expected: "string",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgramForTypes(t, tt.input)
			tc := NewTypeChecker()

			resultType, errors := tc.CheckProgram(program)

			if tt.hasError {
				if len(errors) == 0 {
					t.Errorf("expected errors, but got none")
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("unexpected errors: %v", errors)
				}
			}

			if resultType.String() != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, resultType.String())
			}
		})
	}
}

// TestTypeChecker_ArithmeticOperations は算術演算の型検査をテストする
func TestTypeChecker_ArithmeticOperations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "整数の加算",
			input:    "5 + 3",
			expected: "int",
			hasError: false,
		},
		{
			name:     "浮動小数点の加算",
			input:    "5.5 + 3.2",
			expected: "float",
			hasError: false,
		},
		{
			name:     "整数と浮動小数点の加算（型昇格）",
			input:    "5 + 3.2",
			expected: "float",
			hasError: false,
		},
		{
			name:     "文字列と数値の加算（エラー）",
			input:    "\"hello\" + 5",
			expected: "int",
			hasError: true,
		},
		{
			name:     "整数の減算",
			input:    "10 - 3",
			expected: "int",
			hasError: false,
		},
		{
			name:     "整数の乗算",
			input:    "6 * 7",
			expected: "int",
			hasError: false,
		},
		{
			name:     "整数の除算",
			input:    "15 / 3",
			expected: "int",
			hasError: false,
		},
		{
			name:     "整数の剰余",
			input:    "17 % 5",
			expected: "int",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgramForTypes(t, tt.input)
			tc := NewTypeChecker()

			resultType, errors := tc.CheckProgram(program)

			if tt.hasError {
				if len(errors) == 0 {
					t.Errorf("expected errors, but got none")
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("unexpected errors: %v", errors)
				}
			}

			if !tt.hasError && resultType.String() != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, resultType.String())
			}
		})
	}
}

// TestTypeChecker_ComparisonOperations は比較演算の型検査をテストする
func TestTypeChecker_ComparisonOperations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "整数の等価比較",
			input:    "5 == 5",
			expected: "bool",
			hasError: false,
		},
		{
			name:     "整数の非等価比較",
			input:    "5 != 3",
			expected: "bool",
			hasError: false,
		},
		{
			name:     "整数の大小比較",
			input:    "5 < 10",
			expected: "bool",
			hasError: false,
		},
		{
			name:     "異なる型の比較（エラー）",
			input:    "5 == \"hello\"",
			expected: "bool",
			hasError: true,
		},
		{
			name:     "文字列と数値の大小比較（エラー）",
			input:    "\"hello\" < 5",
			expected: "bool",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgramForTypes(t, tt.input)
			tc := NewTypeChecker()

			resultType, errors := tc.CheckProgram(program)

			if tt.hasError {
				if len(errors) == 0 {
					t.Errorf("expected errors, but got none")
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("unexpected errors: %v", errors)
				}
			}

			if !tt.hasError && resultType.String() != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, resultType.String())
			}
		})
	}
}

// TestTypeChecker_PrefixOperations は前置演算の型検査をテストする
func TestTypeChecker_PrefixOperations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "整数の否定",
			input:    "-42",
			expected: "int",
			hasError: false,
		},
		{
			name:     "浮動小数点の否定",
			input:    "-3.14",
			expected: "float",
			hasError: false,
		},
		{
			name:     "ブール値の論理否定",
			input:    "!true",
			expected: "bool",
			hasError: false,
		},
		{
			name:     "整数の論理否定（エラー）",
			input:    "!42",
			expected: "bool",
			hasError: true,
		},
		{
			name:     "文字列の否定（エラー）",
			input:    "-\"hello\"",
			expected: "string",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgramForTypes(t, tt.input)
			tc := NewTypeChecker()

			resultType, errors := tc.CheckProgram(program)

			if tt.hasError {
				if len(errors) == 0 {
					t.Errorf("expected errors, but got none")
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("unexpected errors: %v", errors)
				}
			}

			if !tt.hasError && resultType.String() != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, resultType.String())
			}
		})
	}
}

// TestTypeChecker_FunctionTypes は関数型の型検査をテストする
func TestTypeChecker_FunctionTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{
			name:     "単純な関数定義",
			input:    "fn(x) { return x; }",
			hasError: false,
		},
		{
			name:     "複数パラメータの関数",
			input:    "fn(x, y) { return x + y; }",
			hasError: false,
		},
		{
			name:     "戻り値なしの関数",
			input:    "fn() { let x = 42; }",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgramForTypes(t, tt.input)
			tc := NewTypeChecker()

			resultType, errors := tc.CheckProgram(program)

			if tt.hasError {
				if len(errors) == 0 {
					t.Errorf("expected errors, but got none")
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("unexpected errors: %v", errors)
				}

				// 関数型であることを確認
				if _, ok := resultType.(*FunctionType); !ok {
					t.Errorf("expected function type, got %T", resultType)
				}
			}
		})
	}
}

// TestTypeChecker_BuiltinFunctions は組み込み関数の型検査をテストする
func TestTypeChecker_BuiltinFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "puts関数の呼び出し",
			input:    "puts(\"hello\")",
			expected: "string",
			hasError: false,
		},
		{
			name:     "puts関数に複数の引数",
			input:    "puts(\"hello\", 42, true)",
			expected: "string",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgramForTypes(t, tt.input)
			tc := NewTypeChecker()

			resultType, errors := tc.CheckProgram(program)

			if tt.hasError {
				if len(errors) == 0 {
					t.Errorf("expected errors, but got none")
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("unexpected errors: %v", errors)
				}
			}

			if !tt.hasError && resultType.String() != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, resultType.String())
			}
		})
	}
}

// TestTypeChecker_IfExpressions はif式の型検査をテストする
func TestTypeChecker_IfExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "条件がブール値のif式",
			input:    "if (true) { 42 }",
			expected: "int",
			hasError: false,
		},
		{
			name:     "条件が整数のif式（エラー）",
			input:    "if (42) { 100 }",
			expected: "int",
			hasError: true,
		},
		{
			name:     "if-else式（同じ型）",
			input:    "if (true) { 42 } else { 24 }",
			expected: "int",
			hasError: false,
		},
		{
			name:     "if-else式（異なる型・エラー）",
			input:    "if (true) { 42 } else { \"hello\" }",
			expected: "string",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgramForTypes(t, tt.input)
			tc := NewTypeChecker()

			resultType, errors := tc.CheckProgram(program)

			if tt.hasError {
				if len(errors) == 0 {
					t.Errorf("expected errors, but got none")
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("unexpected errors: %v", errors)
				}
			}

			if !tt.hasError && resultType.String() != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, resultType.String())
			}
		})
	}
}

// TestTypeChecker_ComplexExpressions は複雑な式の型検査をテストする
func TestTypeChecker_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "複雑な算術式",
			input:    "let x = 10; let y = 20; (x + y) * 2 - 5",
			expected: "int",
			hasError: false,
		},
		{
			name:     "変数と比較演算の組み合わせ",
			input:    "let a = 5; let b = 10; a < b",
			expected: "bool",
			hasError: false,
		},
		{
			name:     "ネストしたif式",
			input:    "if (true) { if (false) { 1 } else { 2 } } else { 3 }",
			expected: "int",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgramForTypes(t, tt.input)
			tc := NewTypeChecker()

			resultType, errors := tc.CheckProgram(program)

			if tt.hasError {
				if len(errors) == 0 {
					t.Errorf("expected errors, but got none")
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("unexpected errors: %v", errors)
				}
			}

			if !tt.hasError && resultType.String() != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, resultType.String())
			}
		})
	}
}

// TestTypeChecker_ErrorMessages は型エラーメッセージをテストする
func TestTypeChecker_ErrorMessages(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedErrors []string
	}{
		{
			name:  "未定義変数",
			input: "undefined_var",
			expectedErrors: []string{
				"identifier not found: undefined_var",
			},
		},
		{
			name:  "型不一致の算術演算",
			input: "\"hello\" + 5",
			expectedErrors: []string{
				"left operand of + must be numeric",
			},
		},
		{
			name:  "型不一致の比較演算",
			input: "5 == \"hello\"",
			expectedErrors: []string{
				"cannot compare int with string",
			},
		},
		{
			name:  "論理否定の型エラー",
			input: "!42",
			expectedErrors: []string{
				"operand of ! must be bool",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgramForTypes(t, tt.input)
			tc := NewTypeChecker()

			_, errors := tc.CheckProgram(program)

			if len(errors) == 0 {
				t.Errorf("expected errors, but got none")
				return
			}

			for _, expectedError := range tt.expectedErrors {
				found := false
				for _, actualError := range errors {
					if strings.Contains(actualError, expectedError) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing '%s', but got errors: %v", expectedError, errors)
				}
			}
		})
	}
}

// TestTypeChecker_TypeEnvironment は型環境のテストを行う
func TestTypeChecker_TypeEnvironment(t *testing.T) {
	env := NewTypeEnvironment()

	// 変数を設定
	env.Set("x", INT_TYPE)
	env.Set("y", STRING_TYPE)

	// 変数を取得
	xType, ok := env.Get("x")
	if !ok {
		t.Errorf("expected to find variable 'x'")
	}
	if !xType.Equals(INT_TYPE) {
		t.Errorf("expected type int for variable 'x', got %s", xType.String())
	}

	// 存在しない変数
	_, ok = env.Get("z")
	if ok {
		t.Errorf("expected not to find variable 'z'")
	}

	// 入れ子の環境
	innerEnv := NewEnclosedTypeEnvironment(env)
	innerEnv.Set("z", BOOL_TYPE)

	// 内側の環境から外側の変数にアクセス
	_, ok = innerEnv.Get("x")
	if !ok {
		t.Errorf("expected to find variable 'x' in enclosed environment")
	}

	// 内側の変数
	zType, ok := innerEnv.Get("z")
	if !ok {
		t.Errorf("expected to find variable 'z' in enclosed environment")
	}
	if !zType.Equals(BOOL_TYPE) {
		t.Errorf("expected type bool for variable 'z', got %s", zType.String())
	}
}

// TestType_Equals は型の等価性テストを行う
func TestType_Equals(t *testing.T) {
	tests := []struct {
		type1    Type
		type2    Type
		expected bool
	}{
		{INT_TYPE, INT_TYPE, true},
		{INT_TYPE, FLOAT_TYPE, false},
		{STRING_TYPE, STRING_TYPE, true},
		{BOOL_TYPE, INT_TYPE, false},
		{
			&FunctionType{Parameters: []Type{INT_TYPE}, ReturnType: STRING_TYPE},
			&FunctionType{Parameters: []Type{INT_TYPE}, ReturnType: STRING_TYPE},
			true,
		},
		{
			&FunctionType{Parameters: []Type{INT_TYPE}, ReturnType: STRING_TYPE},
			&FunctionType{Parameters: []Type{STRING_TYPE}, ReturnType: STRING_TYPE},
			false,
		},
		{
			&ArrayType{ElementType: INT_TYPE},
			&ArrayType{ElementType: INT_TYPE},
			true,
		},
		{
			&ArrayType{ElementType: INT_TYPE},
			&ArrayType{ElementType: STRING_TYPE},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			result := tt.type1.Equals(tt.type2)
			if result != tt.expected {
				t.Errorf("expected %s.Equals(%s) = %v, got %v",
					tt.type1.String(), tt.type2.String(), tt.expected, result)
			}
		})
	}
}

// BenchmarkTypeChecker_BasicTypes は基本型検査のベンチマークを行う
func BenchmarkTypeChecker_BasicTypes(b *testing.B) {
	input := "let x = 42; let y = 3.14; let z = \"hello\"; x + y"
	program := parseProgramForTypes(&testing.T{}, input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := NewTypeChecker()
		_, _ = tc.CheckProgram(program)
	}
}

// BenchmarkTypeChecker_ComplexExpression は複雑な式の型検査ベンチマークを行う
func BenchmarkTypeChecker_ComplexExpression(b *testing.B) {
	input := `
		let a = 10;
		let b = 20;
		let c = 30;
		let result = (a + b) * c - (a * b) / c + (a - b) % c;
		if (result > 100) {
			result - 50
		} else {
			result + 50
		}
	`
	program := parseProgramForTypes(&testing.T{}, input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := NewTypeChecker()
		_, _ = tc.CheckProgram(program)
	}
}
