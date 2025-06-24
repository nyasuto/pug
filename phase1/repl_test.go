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
			hasErrorOutput := strings.Contains(result, "\t") // Error messages are indented

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

	// Check that errors are properly formatted with tabs
	for _, err := range errors {
		expectedError := "\t" + err + "\n"
		if !strings.Contains(result, expectedError) {
			t.Errorf("Expected formatted error %q, got:\n%s", expectedError, result)
		}
	}
}
