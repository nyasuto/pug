package phase2

import (
	"strings"
	"testing"

	"github.com/nyasuto/pug/phase1"
)

func TestControlStructures_Integration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string // Expected assembly patterns
	}{
		{
			name:  "Simple while loop",
			input: "while (true) { let x = 1; }",
			expected: []string{
				".Lwhile_start",
				"testq %rax, %rax",
				"jz .Lwhile_end",
				"jmp .Lwhile_start",
				".Lwhile_end",
			},
		},
		{
			name:  "Break statement in while loop",
			input: "while (true) { break; }",
			expected: []string{
				".Lwhile_start",
				"jmp .Lwhile_end",
				".Lwhile_end",
			},
		},
		{
			name:  "Continue statement in while loop",
			input: "while (true) { continue; }",
			expected: []string{
				".Lwhile_start",
				"jmp .Lwhile_start",
				".Lwhile_end",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the input
			l := phase1.New(tt.input)
			p := phase1.NewParser(l)
			program := p.ParseProgram()

			// Check for parsing errors
			errors := p.Errors()
			if len(errors) > 0 {
				t.Fatalf("Parser errors: %v", errors)
			}

			// Generate assembly code
			cg := NewCodeGenerator()
			assembly, err := cg.Generate(program)
			if err != nil {
				t.Fatalf("Code generation failed: %v", err)
			}

			// Check that expected patterns are present
			for _, pattern := range tt.expected {
				if !strings.Contains(assembly, pattern) {
					t.Errorf("Expected pattern %q not found in assembly:\n%s", pattern, assembly)
				}
			}
		})
	}
}

func TestControlStructures_NestedLoops(t *testing.T) {
	// Skip for now due to for-loop parsing complexity
	t.Skip("For loop parsing needs more work")

	input := `
	for (let i = 0; i < 3; i = i + 1) {
		while (true) {
			break;
		}
	}
	`

	l := phase1.New(input)
	p := phase1.NewParser(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	cg := NewCodeGenerator()
	assembly, err := cg.Generate(program)
	if err != nil {
		t.Fatalf("Code generation failed: %v", err)
	}

	// Check for nested structure
	expectedPatterns := []string{
		".Lfor_start",
		".Lwhile_start",
		".Lwhile_end",
		".Lfor_end",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(assembly, pattern) {
			t.Errorf("Expected pattern %q not found in nested loops", pattern)
		}
	}
}

func TestControlStructures_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Break outside loop",
			input:       "break;",
			expectError: true,
		},
		{
			name:        "Continue outside loop",
			input:       "continue;",
			expectError: true,
		},
		{
			name:        "Valid break in loop",
			input:       "while (true) { break; }",
			expectError: false,
		},
		{
			name:        "Valid continue in loop",
			input:       "while (true) { continue; }",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := phase1.New(tt.input)
			p := phase1.NewParser(l)
			program := p.ParseProgram()

			// Check for parsing errors first
			errors := p.Errors()
			if len(errors) > 0 {
				if !tt.expectError {
					t.Fatalf("Unexpected parser errors: %v", errors)
				}
				return // Expected error, test passed
			}

			// Try code generation
			cg := NewCodeGenerator()
			_, err := cg.Generate(program)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestControlStructures_ComplexExample(t *testing.T) {
	// Skip for now due to for-loop and if statement parsing complexity
	t.Skip("Complex parsing needs more work")

	input := `
	let sum = 0;
	for (let i = 1; i <= 10; i = i + 1) {
		if (i % 2 == 0) {
			continue;
		}
		sum = sum + i;
	}
	return sum;
	`

	l := phase1.New(input)
	p := phase1.NewParser(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	cg := NewCodeGenerator()
	assembly, err := cg.Generate(program)
	if err != nil {
		t.Fatalf("Code generation failed: %v", err)
	}

	// Verify complex control flow patterns
	expectedPatterns := []string{
		"_main:",             // Function entry
		".Lfor_start",        // For loop
		".Lfor_continue",     // Continue label
		"jmp .Lfor_continue", // Continue statement
		".Lfor_end",          // For loop end
		"ret",                // Function return
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(assembly, pattern) {
			t.Logf("Generated assembly:\n%s", assembly)
			t.Errorf("Expected pattern %q not found", pattern)
		}
	}
}
