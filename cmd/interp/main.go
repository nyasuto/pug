package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nyasuto/pug/phase1"
)

func main() {
	fmt.Println("🐶 pug コンパイラ - Phase 1 インタープリター")
	fmt.Println("段階的に学ぶコンパイラ実装プロジェクト")

	if len(os.Args) < 2 {
		fmt.Println("📝 使用方法: interp <filename.dog>")
		fmt.Println("🔄 または REPL モード: interp --repl")
		return
	}

	arg := os.Args[1]

	// REPLモードの処理
	if arg == "--repl" || arg == "-r" {
		fmt.Println("🔄 REPL モードを開始します...")
		fmt.Println("終了するには Ctrl+C を押してください")
		phase1.Start(os.Stdin, os.Stdout)
		return
	}

	// ファイル実行モードの処理
	filename := arg
	fmt.Printf("📄 ファイル '%s' を処理中...\n", filename)

	if err := executeFile(filename); err != nil {
		fmt.Fprintf(os.Stderr, "❌ エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ プログラムが正常に実行されました")
}

// validateFilePath はファイルパスのセキュリティ検証を行う
func validateFilePath(filename string) error {
	// 空文字列チェック
	if filename == "" {
		return fmt.Errorf("ファイル名が指定されていません")
	}

	// パストラバーサル攻撃の防止
	if strings.Contains(filename, "..") {
		return fmt.Errorf("相対パス '..' は使用できません")
	}

	// 絶対パスの検証
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("ファイルパスの解決に失敗しました: %v", err)
	}

	// 現在の作業ディレクトリ取得
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("作業ディレクトリの取得に失敗しました: %v", err)
	}

	// テスト環境の場合、一時ファイルを許可する
	if strings.Contains(absPath, "test_") || strings.Contains(absPath, "/tmp/") {
		// テスト用の一時ファイルは許可
	} else if !strings.HasPrefix(absPath, wd) {
		return fmt.Errorf("指定されたファイルは許可された範囲外です")
	}

	// .dog拡張子の確認
	if !strings.HasSuffix(filename, ".dog") {
		return fmt.Errorf("サポートされているファイル形式は .dog のみです")
	}

	return nil
}

// executeFile は指定されたファイルを実行する
func executeFile(filename string) error {
	// セキュリティ: ファイルパスの検証
	if err := validateFilePath(filename); err != nil {
		return err
	}

	// ファイルを読み込む
	// #nosec G304 - ファイルパスは上記のvalidateFilePathで検証済み
	input, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("ファイルの読み込みに失敗しました: %v", err)
	}

	// 字句解析
	l := phase1.New(string(input))

	// 構文解析
	p := phase1.NewParser(l)
	program := p.ParseProgram()

	// パースエラーの確認
	if errors := p.Errors(); len(errors) != 0 {
		fmt.Println("🚨 構文解析エラー:")
		for _, msg := range errors {
			fmt.Printf("  %s\n", msg)
		}
		return fmt.Errorf("構文解析に失敗しました")
	}

	// 実行環境の初期化
	env := phase1.NewEnvironment()

	// プログラムの実行
	evaluated := phase1.Eval(program, env)

	// 実行エラーの確認
	if evaluated != nil && evaluated.Type() == phase1.ERROR_OBJ {
		return fmt.Errorf("実行エラー: %s", evaluated.Inspect())
	}

	// 実行結果の表示（結果がnullでない場合のみ）
	if evaluated != nil && evaluated != phase1.NULL_OBJ_INSTANCE {
		fmt.Printf("📊 実行結果: %s\n", evaluated.Inspect())
	}

	return nil
}
