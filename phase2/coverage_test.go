package phase2

import (
	"strings"
	"testing"

	"github.com/nyasuto/pug/phase1"
)

// TestCodeGeneratorCoverage tests various code paths to improve coverage
func TestCodeGeneratorCoverage(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Boolean true",
			input: "true",
		},
		{
			name:  "Boolean false",
			input: "false",
		},
		{
			name:  "Modulo operation",
			input: "10 % 3",
		},
		{
			name:  "Division operation",
			input: "10 / 2",
		},
		{
			name:  "Logical negation",
			input: "!true",
		},
		{
			name:  "Arithmetic negation",
			input: "-42",
		},
		{
			name:  "Not equal comparison",
			input: "5 != 3",
		},
		{
			name:  "Less than comparison",
			input: "5 < 10",
		},
		{
			name:  "Greater than comparison",
			input: "10 > 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := phase1.New(tt.input)
			p := phase1.NewParser(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			cg := NewCodeGenerator()
			_, err := cg.Generate(program)
			if err != nil {
				t.Fatalf("Code generation failed: %v", err)
			}
		})
	}
}

// TestCodeGeneratorErrorPaths tests error handling paths
func TestCodeGeneratorErrorPaths(t *testing.T) {
	cg := NewCodeGenerator()

	// Test undefined variable
	ident := &phase1.Identifier{
		Token: phase1.Token{Type: phase1.IDENT, Literal: "undefined"},
		Value: "undefined",
	}
	err := cg.generateIdentifier(ident)
	if err == nil {
		t.Errorf("Expected error for undefined variable")
	}
	if !strings.Contains(err.Error(), "undefined variable") {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

// TestControlFlowValidation tests control flow validation
func TestControlFlowValidation(t *testing.T) {
	analyzer := NewControlFlowAnalyzer()

	// Test nested scopes
	analyzer.EnterScope()
	analyzer.EnterScope()
	if analyzer.GetSymbolTable().GetScopeLevel() != 2 {
		t.Errorf("Expected scope level 2, got %d", analyzer.GetSymbolTable().GetScopeLevel())
	}

	analyzer.ExitScope()
	analyzer.ExitScope()
	if analyzer.GetSymbolTable().GetScopeLevel() != 0 {
		t.Errorf("Expected scope level 0 after exiting all scopes")
	}

	// Test nested loops
	analyzer.EnterLoop("outer_break", "outer_continue")
	analyzer.EnterLoop("inner_break", "inner_continue")

	ctx := analyzer.GetCurrentLoopContext()
	if ctx.BreakLabel != "inner_break" {
		t.Errorf("Expected inner loop context")
	}

	analyzer.ExitLoop()
	ctx = analyzer.GetCurrentLoopContext()
	if ctx.BreakLabel != "outer_break" {
		t.Errorf("Expected outer loop context after exiting inner")
	}

	analyzer.ExitLoop()
	if analyzer.GetCurrentLoopContext() != nil {
		t.Errorf("Expected no loop context after exiting all loops")
	}
}
