package phase1

import (
	"bytes"
	"strings"
	"testing"
)

func TestREPLStart(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple arithmetic",
			input:    "5 + 5\n",
			expected: []string{"10"},
		},
		{
			name:     "variable assignment",
			input:    "let x = 10\nx\n",
			expected: []string{"10", "10"},
		},
		{
			name:     "function definition and call",
			input:    "let add = fn(a, b) { a + b }\nadd(3, 4)\n",
			expected: []string{"fn(a, b) {\n(a + b)\n}", "7"},
		},
		{
			name:     "string operations",
			input:    "\"hello\" + \" \" + \"world\"\n",
			expected: []string{"hello world"},
		},
		{
			name:     "boolean operations",
			input:    "true && false\n!true\n",
			expected: []string{"false", "false"},
		},
		{
			name:     "built-in functions",
			input:    "len(\"hello\")\ntype(42)\n",
			expected: []string{"5", "INTEGER"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			output := &bytes.Buffer{}

			// Capture the REPL output
			Start(input, output)

			result := output.String()

			// Check that expected outputs are present
			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected output to contain %q, got:\n%s", expected, result)
				}
			}
		})
	}
}

func TestREPLParserErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{
			name:     "syntax error",
			input:    "let x =\n",
			hasError: true,
		},
		{
			name:     "incomplete expression",
			input:    "5 +\n",
			hasError: true,
		},
		{
			name:     "valid expression",
			input:    "5 + 5\n",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			output := &bytes.Buffer{}

			Start(input, output)

			result := output.String()
			hasErrorOutput := strings.Contains(result, "🚨 構文解析エラー:") // Error messages start with this

			if tt.hasError && !hasErrorOutput {
				t.Errorf("Expected error output, but got none. Output:\n%s", result)
			}
			if !tt.hasError && hasErrorOutput {
				t.Errorf("Expected no error output, but got errors. Output:\n%s", result)
			}
		})
	}
}

func TestREPLPersistentEnvironment(t *testing.T) {
	// Test that variables persist across REPL lines
	input := strings.NewReader("let x = 42\nlet y = x + 8\ny\n")
	output := &bytes.Buffer{}

	Start(input, output)

	result := output.String()

	// Should contain the final result of 50
	if !strings.Contains(result, "50") {
		t.Errorf("Expected persistent environment to work, got:\n%s", result)
	}
}

func TestREPLMultipleStatements(t *testing.T) {
	// Test multiple statements in one session
	input := strings.NewReader("5\n10\n15\n")
	output := &bytes.Buffer{}

	Start(input, output)

	result := output.String()

	// Check that all results are present
	expectedResults := []string{"5", "10", "15"}
	for _, expected := range expectedResults {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected output to contain %q, got:\n%s", expected, result)
		}
	}
}

func TestPrintParserErrors(t *testing.T) {
	output := &bytes.Buffer{}
	errors := []string{
		"expected next token to be =, got INT instead",
		"no prefix parse function for INT found",
	}

	printParserErrors(output, errors)

	result := output.String()

	// Check that errors are properly formatted
	for _, err := range errors {
		if !strings.Contains(result, err) {
			t.Errorf("Expected error %q, got:\n%s", err, result)
		}
	}
}

// Enhanced REPL tests for new functionality

func TestREPLSpecialCommands(t *testing.T) {
	tests := []struct {
		command  string
		expected string
	}{
		{":help", "📖 Pug REPL ヘルプ:"},
		{":exit", "👋 Goodbye!"},
		{":quit", "👋 Goodbye!"},
		{":q", "👋 Goodbye!"},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			in := strings.NewReader(tt.command + "\n")
			var out bytes.Buffer

			Start(in, &out)

			output := out.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestREPLHistory(t *testing.T) {
	input := "5 + 3\nlet x = 42\n:history\n:exit\n"
	in := strings.NewReader(input)
	var out bytes.Buffer

	Start(in, &out)

	output := out.String()

	// 履歴に両方のコマンドが含まれることを確認
	if !strings.Contains(output, "1: 5 + 3") {
		t.Error("Expected history to contain '1: 5 + 3'")
	}
	if !strings.Contains(output, "2: let x = 42") {
		t.Error("Expected history to contain '2: let x = 42'")
	}
}

func TestREPLEnhancedErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"unknownVariable", "identifier not found: unknownVariable"},
		{"5 / 0", "division by zero"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			in := strings.NewReader(tt.input + "\n:exit\n")
			var out bytes.Buffer

			Start(in, &out)

			output := out.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestREPLUnknownCommand(t *testing.T) {
	input := ":unknown\n:exit\n"
	in := strings.NewReader(input)
	var out bytes.Buffer

	Start(in, &out)

	output := out.String()

	// 不明なコマンドに対する適切なエラーメッセージ
	if !strings.Contains(output, "❓ 不明なコマンド: :unknown") {
		t.Error("Expected unknown command error message")
	}
}

func TestREPLEmptyLines(t *testing.T) {
	input := "\n\n5 + 3\n\n:exit\n"
	in := strings.NewReader(input)
	var out bytes.Buffer

	Start(in, &out)

	output := out.String()

	// 空行がスキップされ、計算が正しく実行されることを確認
	if !strings.Contains(output, "=> 8") {
		t.Error("Expected 5 + 3 to equal 8 despite empty lines")
	}
}

func TestREPLEnhancedPromptAndOutput(t *testing.T) {
	input := "42\n:exit\n"
	in := strings.NewReader(input)
	var out bytes.Buffer

	Start(in, &out)

	output := out.String()

	// 新しいプロンプトとアウトプット形式の確認
	if !strings.Contains(output, "pug>") {
		t.Error("Expected new prompt 'pug>'")
	}
	if !strings.Contains(output, "=> 42") {
		t.Error("Expected enhanced output format '=> 42'")
	}
}
