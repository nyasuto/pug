package main

import (
	"fmt"
	"os"

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

// executeFile は指定されたファイルを実行する
func executeFile(filename string) error {
	// ファイルを読み込む
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
