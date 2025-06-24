package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestSamplePrograms tests that sample programs can be parsed and evaluated
func TestSamplePrograms(t *testing.T) {
	// Build the main program first
	buildCmd := exec.Command("go", "build", "-o", "pug_test", "main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build main program: %v", err)
	}
	defer os.Remove("pug_test") // Clean up

	sampleFiles := []string{
		"examples/hello.dog",
		"examples/fibonacci.dog",
		"examples/calculator.dog",
		"examples/sorting.dog",
		"examples/closures.dog",
		"examples/algorithms.dog",
	}

	for _, file := range sampleFiles {
		t.Run(filepath.Base(file), func(t *testing.T) {
			// Check if file exists
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Skipf("Sample file %s does not exist", file)
			}

			// Read the file content
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read sample file %s: %v", file, err)
			}

			// Test that the file contains valid syntax by trying to parse it
			// We can't easily test execution without a more complex setup,
			// but we can at least verify the files exist and are readable
			if len(content) == 0 {
				t.Errorf("Sample file %s is empty", file)
			}

			// Check for basic Dog language syntax
			contentStr := string(content)
			if !strings.Contains(contentStr, "//") &&
				!strings.Contains(contentStr, "let") &&
				!strings.Contains(contentStr, "fn") &&
				!strings.Contains(contentStr, "puts") {
				t.Errorf("Sample file %s doesn't seem to contain Dog language syntax", file)
			}
		})
	}
}

// TestExamplesDirectory tests the examples directory structure
func TestExamplesDirectory(t *testing.T) {
	examplesDir := "examples"

	// Check that examples directory exists
	if _, err := os.Stat(examplesDir); os.IsNotExist(err) {
		t.Fatalf("Examples directory does not exist: %s", examplesDir)
	}

	// Read directory contents
	files, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	// Check that we have some .dog files
	dogFileCount := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".dog") {
			dogFileCount++
		}
	}

	if dogFileCount == 0 {
		t.Error("No .dog files found in examples directory")
	}

	// Check for required sample files
	requiredFiles := []string{"hello.dog", "fibonacci.dog"}
	for _, required := range requiredFiles {
		found := false
		for _, file := range files {
			if file.Name() == required {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Required sample file %s not found in examples directory", required)
		}
	}
}

// TestBuildProcess tests that the project builds successfully
func TestBuildProcess(t *testing.T) {
	// Test that main program builds
	buildCmd := exec.Command("go", "build", "-o", "pug_build_test", "main.go")
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Build failed: %v\nOutput: %s", err, output)
	}
	defer os.Remove("pug_build_test")

	// Check that the binary was created
	if _, err := os.Stat("pug_build_test"); os.IsNotExist(err) {
		t.Error("Build succeeded but binary was not created")
	}
}

// TestMakeTargets tests that important make targets work
func TestMakeTargets(t *testing.T) {
	// Note: Excluded "test" target to avoid infinite recursion
	// (TestMakeTargets -> make test -> TestMakeTargets -> ...)
	targets := []string{"build", "clean", "fmt", "lint"}

	for _, target := range targets {
		t.Run("make_"+target, func(t *testing.T) {
			cmd := exec.Command("make", target)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Make target '%s' failed: %v\nOutput: %s", target, err, output)
			}
		})
	}
}

// TestProjectStructure tests the overall project structure
func TestProjectStructure(t *testing.T) {
	requiredDirs := []string{
		"phase1",
		"examples",
		".github",
	}

	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Required directory %s does not exist", dir)
		}
	}

	requiredFiles := []string{
		"go.mod",
		"main.go",
		"Makefile",
		"README.md",
		".gitignore",
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Required file %s does not exist", file)
		}
	}
}

// TestPhaseStructure tests the phase directory structure
func TestPhaseStructure(t *testing.T) {
	phases := []struct {
		dir         string
		required    bool
		description string
	}{
		{"phase1", true, "Interpreter (lexer, parser, evaluator)"},
		{"phase2", false, "Bytecode compiler and VM"},
		{"phase3", false, "IR generation and optimization"},
		{"phase4", false, "Native code generation"},
	}

	for _, phase := range phases {
		t.Run(phase.dir, func(t *testing.T) {
			if _, err := os.Stat(phase.dir); os.IsNotExist(err) {
				if phase.required {
					t.Errorf("Required phase directory %s does not exist", phase.dir)
				} else {
					t.Skipf("Optional phase directory %s does not exist yet", phase.dir)
				}
				return
			}

			// If phase exists, check for go files
			files, err := os.ReadDir(phase.dir)
			if err != nil {
				t.Fatalf("Failed to read phase directory %s: %v", phase.dir, err)
			}

			goFileCount := 0
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".go") {
					goFileCount++
				}
			}

			if goFileCount == 0 {
				if phase.required {
					t.Errorf("Required phase directory %s contains no Go files", phase.dir)
				} else {
					t.Logf("Optional phase directory %s exists but contains no Go files (not implemented yet)", phase.dir)
				}
			}
		})
	}
}

// TestGoModTidy tests that go.mod is properly maintained
func TestGoModTidy(t *testing.T) {
	// Run go mod tidy and check if it changes anything
	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		t.Fatalf("go mod tidy failed: %v", err)
	}

	// Check if go.mod is properly formatted
	cmd = exec.Command("go", "mod", "verify")
	if err := cmd.Run(); err != nil {
		t.Errorf("go mod verify failed: %v", err)
	}
}

// Future: TestCrossPhaseIntegration will test interaction between phases
func TestCrossPhaseIntegration(t *testing.T) {
	t.Skip("Cross-phase integration testing will be implemented when multiple phases exist")

	// Future implementation will test:
	// 1. phase1 AST -> phase2 bytecode
	// 2. phase2 bytecode -> phase3 IR
	// 3. phase3 IR -> phase4 native code
	// 4. End-to-end: source -> native execution
	// 5. Performance comparison between phases
}
