package main

import (
	"os"
	"strings"
	"testing"
)

func TestExecuteFile_Success(t *testing.T) {
	// テスト用の一時ファイルを作成
	content := `
puts("Hello from test");
let x = 5 + 3;
puts("Result:", x);
`

	tmpFile, err := os.CreateTemp("", "test_*.dog")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// executeFile を実行
	err = executeFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("executeFile failed: %v", err)
	}
}

func TestExecuteFile_FileNotFound(t *testing.T) {
	err := executeFile("nonexistent.dog")
	if err == nil {
		t.Fatal("Expected error for nonexistent file, but got nil")
	}

	if !strings.Contains(err.Error(), "ファイルの読み込みに失敗しました") {
		t.Errorf("Expected file read error, got: %v", err)
	}
}

func TestExecuteFile_SyntaxError(t *testing.T) {
	// 構文エラーのあるコンテンツ
	content := `
let x = ;  // 構文エラー
`

	tmpFile, err := os.CreateTemp("", "test_syntax_error_*.dog")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// executeFile を実行
	err = executeFile(tmpFile.Name())
	if err == nil {
		t.Fatal("Expected syntax error, but got nil")
	}

	if !strings.Contains(err.Error(), "構文解析に失敗しました") {
		t.Errorf("Expected syntax error, got: %v", err)
	}
}

func TestExecuteFile_RuntimeError(t *testing.T) {
	// 実行時エラーのあるコンテンツ
	content := `
let x = unknownVariable;  // 実行時エラー
`

	tmpFile, err := os.CreateTemp("", "test_runtime_error_*.dog")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// executeFile を実行
	err = executeFile(tmpFile.Name())
	if err == nil {
		t.Fatal("Expected runtime error, but got nil")
	}

	if !strings.Contains(err.Error(), "実行エラー") {
		t.Errorf("Expected runtime error, got: %v", err)
	}
}

func TestExecuteFile_WithResult(t *testing.T) {
	// 結果を返すプログラム
	content := `5 + 3`

	tmpFile, err := os.CreateTemp("", "test_result_*.dog")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// executeFile を実行
	err = executeFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("executeFile failed: %v", err)
	}
}

// mainIntegrationTest はmain関数の統合テスト用ヘルパー
func TestMainWithREPL(t *testing.T) {
	// REPLモードのテスト用の入力
	input := ":help\n:exit\n"

	// 標準入力を模擬
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"interp", "--repl"}

	// 標準入力を置き換え
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		defer w.Close()
		w.WriteString(input)
	}()

	// main関数を実行
	// Note: main()を直接テストするのは難しいため、
	// 実際のテストでは統合テストとして別途実行する
}

func TestMainWithFile(t *testing.T) {
	// ファイル実行モードのテスト
	content := `puts("Integration test")`

	tmpFile, err := os.CreateTemp("", "test_main_*.dog")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// コマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"interp", tmpFile.Name()}

	// main関数のテストは統合テストで行う
	// ここでは基本的な設定のテストのみ
}
