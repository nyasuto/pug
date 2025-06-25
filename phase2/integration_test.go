package phase2

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestIntegration_EndToEnd はエンドツーエンドの統合テストを実行する
func TestIntegration_EndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// CI環境では統合テストをスキップ
	if os.Getenv("CI") != "" {
		t.Skip("skipping assembly integration tests in CI environment (assembly generation tests provide sufficient coverage)")
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple_return",
			input:    "return 42;",
			expected: "42",
		},
		{
			name:     "arithmetic",
			input:    "return 5 + 3 * 2;",
			expected: "11",
		},
		{
			name:     "variables",
			input:    "let x = 10; let y = 20; return x + y;",
			expected: "30",
		},
		{
			name:     "comparison",
			input:    "return 5 == 5;",
			expected: "1",
		},
		{
			name:     "complex_expression",
			input:    "let a = 5; let b = 3; return (a + b) * 2 - 1;",
			expected: "15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テンポラリディレクトリを作成
			tmpDir, err := os.MkdirTemp("", "pug_test_"+tt.name)
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// アセンブリコードを生成
			program := parseProgram(t, tt.input)
			cg := NewCodeGenerator()
			asmCode, err := cg.Generate(program)
			if err != nil {
				t.Fatalf("code generation failed: %v", err)
			}

			// アセンブリファイルに書き込み
			asmFile := filepath.Join(tmpDir, tt.name+".s")
			err = os.WriteFile(asmFile, []byte(asmCode), 0644)
			if err != nil {
				t.Fatalf("failed to write assembly file: %v", err)
			}

			// アセンブルとリンクを実行
			execFile := filepath.Join(tmpDir, tt.name)
			err = assembleAndLink(asmFile, execFile)
			if err != nil {
				t.Logf("Assembly code:\n%s", asmCode)
				t.Fatalf("failed to assemble and link: %v", err)
			}

			// 実行可能ファイルを実行
			output, exitCode, err := runExecutable(execFile)
			if err != nil {
				t.Fatalf("failed to run executable: %v", err)
			}

			// 戻り値を確認（終了コードとして返される）
			expectedCode := tt.expected
			actualCode := fmt.Sprintf("%d", exitCode)

			if actualCode != expectedCode {
				t.Errorf("expected exit code %s, but got %s. Output: %s", expectedCode, actualCode, output)
			}
		})
	}
}

// TestIntegration_AssemblyGeneration はアセンブリ生成の詳細をテストする
func TestIntegration_AssemblyGeneration(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:  "basic_structure",
			input: "return 1;",
			shouldContain: []string{
				"# pug compiler generated assembly",
				".section __DATA,__data",
				".section __TEXT,__text,regular,pure_instructions",
				".globl _main",
				"_main:",
				"pushq %rbp",
				"movq %rsp, %rbp",
				"subq $256, %rsp",
				"movq $1, %rax",
				"movq %rbp, %rsp",
				"popq %rbp",
				"ret",
			},
			shouldNotContain: []string{
				"undefined",
				"error",
			},
		},
		{
			name:  "variable_management",
			input: "let x = 42; let y = 24; return x + y;",
			shouldContain: []string{
				"movq $42, %rax",
				"movq %rax, -8(%rbp)",
				"# let x = ...",
				"movq $24, %rax",
				"movq %rax, -16(%rbp)",
				"# let y = ...",
				"movq -8(%rbp), %rax",
				"# load variable x",
				"movq -16(%rbp), %rax",
				"# load variable y",
				"addq %rbx, %rax",
			},
		},
		{
			name:  "arithmetic_operations",
			input: "return 10 * 5 / 2;",
			shouldContain: []string{
				"imulq %rbx, %rax",
				"cqto",
				"idivq %rbx",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			cg := NewCodeGenerator()
			asmCode, err := cg.Generate(program)
			if err != nil {
				t.Fatalf("code generation failed: %v", err)
			}

			for _, shouldContain := range tt.shouldContain {
				if !strings.Contains(asmCode, shouldContain) {
					t.Errorf("expected assembly to contain '%s', but it didn't.\nGenerated assembly:\n%s", shouldContain, asmCode)
				}
			}

			for _, shouldNotContain := range tt.shouldNotContain {
				if strings.Contains(asmCode, shouldNotContain) {
					t.Errorf("expected assembly to NOT contain '%s', but it did.\nGenerated assembly:\n%s", shouldNotContain, asmCode)
				}
			}
		})
	}
}

// TestIntegration_Performance は性能比較テストを実行する
func TestIntegration_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	// 複雑な計算のテストプログラム
	input := `
		let a = 10;
		let b = 20;
		let c = 30;
		let result = (a + b) * c - (a * b) / c;
		return result;
	`

	// まずは正常にコンパイルできることを確認
	program := parseProgram(t, input)
	cg := NewCodeGenerator()
	asmCode, err := cg.Generate(program)
	if err != nil {
		t.Fatalf("code generation failed: %v", err)
	}

	// アセンブリコードの長さを確認（効率性の指標）
	lines := strings.Split(asmCode, "\n")
	nonEmptyLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "#") {
			nonEmptyLines++
		}
	}

	t.Logf("Generated %d lines of assembly code", nonEmptyLines)

	// 基本的な効率性チェック（あまりにも長すぎるコードは問題）
	if nonEmptyLines > 200 {
		t.Errorf("generated assembly is too long (%d lines), may be inefficient", nonEmptyLines)
	}
}

// assembleAndLink はアセンブリファイルをアセンブルしてリンクする
func assembleAndLink(asmFile, outputFile string) error {
	// Linux/GNU環境では統合テストをスキップ
	// アセンブリ生成テストのみで十分な検証
	if os.Getenv("CI") != "" {
		return fmt.Errorf("assembly integration tests are skipped in CI environment (assembly generation tests cover the core functionality)")
	}

	// macOSでの実装
	// まずアセンブルしてオブジェクトファイルを作成
	objFile := strings.TrimSuffix(outputFile, filepath.Ext(outputFile)) + ".o"

	// macOSでは as コマンドに -64 フラグは不要
	asCmd := exec.Command("as", "-arch", "x86_64", "-o", objFile, asmFile)
	asOutput, err := asCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("assembly failed: %v\nOutput: %s", err, string(asOutput))
	}

	// リンクして実行可能ファイルを作成
	ldCmd := exec.Command("ld", "-o", outputFile, "-lSystem", "-syslibroot", "/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk", "-e", "_main", "-arch", "x86_64", objFile)
	ldOutput, err := ldCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("linking failed: %v\nOutput: %s", err, string(ldOutput))
	}

	return nil
}

// runExecutable は実行可能ファイルを実行して結果を取得する
func runExecutable(execFile string) (output string, exitCode int, err error) {
	cmd := exec.Command(execFile)
	outputBytes, err := cmd.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
			return string(outputBytes), exitCode, nil
		}
		return "", -1, err
	}

	return string(outputBytes), 0, nil
}

// BenchmarkIntegration_CodeGeneration はコード生成のベンチマークを実行する
func BenchmarkIntegration_CodeGeneration(b *testing.B) {
	// 基本的な式のベンチマーク
	simpleInput := `
		let a = 10;
		let b = 20;
		let c = 30;
		return (a + b) * c - (a * b) / c;
	`

	program := parseProgram(&testing.T{}, simpleInput)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cg := NewCodeGenerator()
		_, err := cg.Generate(program)
		if err != nil {
			b.Fatalf("code generation failed: %v", err)
		}
	}
}

// TestIntegration_ErrorHandling はエラーハンドリングをテストする
func TestIntegration_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "undefined_variable",
			input:       "return undefined_var;",
			expectError: true,
			errorMsg:    "undefined variable",
		},
		{
			name:        "valid_program",
			input:       "let x = 42; return x;",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			cg := NewCodeGenerator()
			_, err := cg.Generate(program)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing '%s', but got none", tt.errorMsg)
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing '%s', but got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, but got: %v", err)
				}
			}
		})
	}
}
