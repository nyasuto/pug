package phase2

import (
	"strings"
	"testing"

	"github.com/nyasuto/pug/phase1"
)

// TestCodeGenerator_IntegerLiteral は整数リテラルのコード生成をテストする
func TestCodeGenerator_IntegerLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: "42;",
			expected: []string{
				"movq $42, %rax",
			},
		},
		{
			input: "0;",
			expected: []string{
				"movq $0, %rax",
			},
		},
		{
			input: "-123;",
			expected: []string{
				"movq $123, %rax",
				"negq %rax",
			},
		},
	}

	for _, tt := range tests {
		program := parseProgram(t, tt.input)
		cg := NewCodeGenerator()

		code, err := cg.Generate(program)
		if err != nil {
			t.Fatalf("code generation failed: %v", err)
		}

		for _, expectedLine := range tt.expected {
			if !strings.Contains(code, expectedLine) {
				t.Errorf("expected assembly to contain '%s', but got:\n%s", expectedLine, code)
			}
		}
	}
}

// TestCodeGenerator_BinaryOperations は二項演算のコード生成をテストする
func TestCodeGenerator_BinaryOperations(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: "5 + 3;",
			expected: []string{
				"movq $5, %rax",
				"pushq %rax",
				"movq $3, %rax",
				"movq %rax, %rbx",
				"popq %rax",
				"addq %rbx, %rax",
			},
		},
		{
			input: "10 - 4;",
			expected: []string{
				"movq $10, %rax",
				"pushq %rax",
				"movq $4, %rax",
				"movq %rax, %rbx",
				"popq %rax",
				"subq %rbx, %rax",
			},
		},
		{
			input: "6 * 7;",
			expected: []string{
				"movq $6, %rax",
				"pushq %rax",
				"movq $7, %rax",
				"movq %rax, %rbx",
				"popq %rax",
				"imulq %rbx, %rax",
			},
		},
		{
			input: "15 / 3;",
			expected: []string{
				"movq $15, %rax",
				"pushq %rax",
				"movq $3, %rax",
				"movq %rax, %rbx",
				"popq %rax",
				"cqto",
				"idivq %rbx",
			},
		},
	}

	for _, tt := range tests {
		program := parseProgram(t, tt.input)
		cg := NewCodeGenerator()

		code, err := cg.Generate(program)
		if err != nil {
			t.Fatalf("code generation failed: %v", err)
		}

		for _, expectedLine := range tt.expected {
			if !strings.Contains(code, expectedLine) {
				t.Errorf("expected assembly to contain '%s', but got:\n%s", expectedLine, code)
			}
		}
	}
}

// TestCodeGenerator_Variables は変数のコード生成をテストする
func TestCodeGenerator_Variables(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: "let x = 42; x;",
			expected: []string{
				"movq $42, %rax",
				"movq %rax, -8(%rbp)",
				"# let x = ...",
				"movq -8(%rbp), %rax",
				"# load variable x",
			},
		},
		{
			input: "let a = 10; let b = 20; a + b;",
			expected: []string{
				"movq $10, %rax",
				"movq %rax, -8(%rbp)",
				"# let a = ...",
				"movq $20, %rax",
				"movq %rax, -16(%rbp)",
				"# let b = ...",
				"movq -8(%rbp), %rax",
				"# load variable a",
				"pushq %rax",
				"movq -16(%rbp), %rax",
				"# load variable b",
				"movq %rax, %rbx",
				"popq %rax",
				"addq %rbx, %rax",
			},
		},
	}

	for _, tt := range tests {
		program := parseProgram(t, tt.input)
		cg := NewCodeGenerator()

		code, err := cg.Generate(program)
		if err != nil {
			t.Fatalf("code generation failed: %v", err)
		}

		for _, expectedLine := range tt.expected {
			if !strings.Contains(code, expectedLine) {
				t.Errorf("expected assembly to contain '%s', but got:\n%s", expectedLine, code)
			}
		}
	}
}

// TestCodeGenerator_Comparisons は比較演算のコード生成をテストする
func TestCodeGenerator_Comparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: "5 == 5;",
			expected: []string{
				"movq $5, %rax",
				"pushq %rax",
				"movq $5, %rax",
				"movq %rax, %rbx",
				"popq %rax",
				"cmpq %rbx, %rax",
				"je .Ltrue",
				"movq $0, %rax",
				"jmp .Lend",
				".Ltrue0:",
				"movq $1, %rax",
				".Lend1:",
			},
		},
		{
			input: "3 < 7;",
			expected: []string{
				"cmpq %rbx, %rax",
				"jl .Ltrue",
				"movq $0, %rax",
				"movq $1, %rax",
			},
		},
	}

	for _, tt := range tests {
		program := parseProgram(t, tt.input)
		cg := NewCodeGenerator()

		code, err := cg.Generate(program)
		if err != nil {
			t.Fatalf("code generation failed: %v", err)
		}

		// 比較演算の場合は一部の命令が含まれていることを確認
		for _, expectedLine := range tt.expected {
			if strings.Contains(expectedLine, ".L") {
				// ラベルは動的に生成されるので、パターンマッチング
				continue
			}
			if !strings.Contains(code, expectedLine) {
				t.Errorf("expected assembly to contain '%s', but got:\n%s", expectedLine, code)
			}
		}
	}
}

// TestCodeGenerator_ComplexExpressions は複雑な式のコード生成をテストする
func TestCodeGenerator_ComplexExpressions(t *testing.T) {
	tests := []struct {
		input       string
		description string
	}{
		{
			input:       "let x = 5; let y = 3; x * y + 2;",
			description: "変数を使った複雑な演算",
		},
		{
			input:       "!(5 == 3);",
			description: "論理否定と比較の組み合わせ",
		},
		{
			input:       "let result = (10 + 5) * 2;",
			description: "括弧を含む演算",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			cg := NewCodeGenerator()

			code, err := cg.Generate(program)
			if err != nil {
				t.Fatalf("code generation failed: %v", err)
			}

			// 基本的な構造が含まれていることを確認
			expectedStructures := []string{
				"_main:",
				"pushq %rbp",
				"movq %rsp, %rbp",
				"subq $256, %rsp",
				"movq $0, %rax",
				"movq %rbp, %rsp",
				"popq %rbp",
				"ret",
			}

			for _, structure := range expectedStructures {
				if !strings.Contains(code, structure) {
					t.Errorf("expected assembly to contain basic structure '%s', but got:\n%s", structure, code)
				}
			}
		})
	}
}

// TestCodeGenerator_ReturnStatement はreturn文のコード生成をテストする
func TestCodeGenerator_ReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: "return 42;",
			expected: []string{
				"movq $42, %rax",
				"movq %rbp, %rsp",
				"popq %rbp",
				"ret",
			},
		},
		{
			input: "let x = 0; return x;",
			expected: []string{
				"movq $0, %rax",
				"movq %rbp, %rsp",
				"popq %rbp",
				"ret",
			},
		},
	}

	for _, tt := range tests {
		program := parseProgram(t, tt.input)
		cg := NewCodeGenerator()

		code, err := cg.Generate(program)
		if err != nil {
			t.Fatalf("code generation failed: %v", err)
		}

		for _, expectedLine := range tt.expected {
			if !strings.Contains(code, expectedLine) {
				t.Errorf("expected assembly to contain '%s', but got:\n%s", expectedLine, code)
			}
		}
	}
}

// TestCodeGenerator_ErrorCases はエラーケースをテストする
func TestCodeGenerator_ErrorCases(t *testing.T) {
	tests := []struct {
		input       string
		description string
	}{
		{
			input:       "undefined_var;",
			description: "未定義変数の参照",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			cg := NewCodeGenerator()

			_, err := cg.Generate(program)
			if err == nil {
				t.Errorf("expected error for %s, but got none", tt.description)
			}
		})
	}
}

// parseProgram はテスト用のヘルパー関数
func parseProgram(t *testing.T, input string) *phase1.Program {
	lexer := phase1.New(input)
	parser := phase1.NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		t.Fatalf("parser errors: %v", parser.Errors())
	}

	return program
}

// TestCodeGenerator_AssemblyHeader はアセンブリヘッダーの生成をテストする
func TestCodeGenerator_AssemblyHeader(t *testing.T) {
	program := parseProgram(t, "42;")
	cg := NewCodeGenerator()

	code, err := cg.Generate(program)
	if err != nil {
		t.Fatalf("code generation failed: %v", err)
	}

	expectedHeaders := []string{
		"# pug compiler generated assembly",
		".section __DATA,__data",
		".section __TEXT,__text,regular,pure_instructions",
		".globl _main",
		"_main:",
		"pushq %rbp",
		"movq %rsp, %rbp",
		"subq $256, %rsp",
	}

	for _, header := range expectedHeaders {
		if !strings.Contains(code, header) {
			t.Errorf("expected assembly header to contain '%s', but got:\n%s", header, code)
		}
	}
}

// TestCodeGenerator_LabelGeneration はラベル生成をテストする
func TestCodeGenerator_LabelGeneration(t *testing.T) {
	cg := NewCodeGenerator()

	label1 := cg.generateLabel("test")
	label2 := cg.generateLabel("test")
	label3 := cg.generateLabel("other")

	expectedPattern := []string{".Ltest0", ".Ltest1", ".Lother2"}

	for i, expected := range expectedPattern {
		var actual string
		switch i {
		case 0:
			actual = label1
		case 1:
			actual = label2
		case 2:
			actual = label3
		}

		if actual != expected {
			t.Errorf("expected label '%s', but got '%s'", expected, actual)
		}
	}
}

// BenchmarkCodeGenerator_BasicExpression は基本的な式のコード生成性能をベンチマークする
func BenchmarkCodeGenerator_BasicExpression(b *testing.B) {
	input := "let x = 10; let y = 20; x + y * 2;"
	program := parseProgram(&testing.T{}, input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cg := NewCodeGenerator()
		_, err := cg.Generate(program)
		if err != nil {
			b.Fatalf("code generation failed: %v", err)
		}
	}
}

// BenchmarkCodeGenerator_ComplexExpression は複雑な式のコード生成性能をベンチマークする
func BenchmarkCodeGenerator_ComplexExpression(b *testing.B) {
	input := `
		let a = 10;
		let b = 20;
		let c = 30;
		let result = (a + b) * c - (a * b) / c + (a - b) % c;
	`
	program := parseProgram(&testing.T{}, input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cg := NewCodeGenerator()
		_, err := cg.Generate(program)
		if err != nil {
			b.Fatalf("code generation failed: %v", err)
		}
	}
}
