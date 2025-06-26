package main

import (
	"fmt"
	"os"

	"github.com/nyasuto/pug/phase1"
	"github.com/nyasuto/pug/phase2"
)

func main() {
	fmt.Println("🐶 pug コンパイラ - Phase 2 コンパイラ")
	fmt.Println("段階的に学ぶコンパイラ実装プロジェクト")

	if len(os.Args) < 2 {
		fmt.Println("📝 使用方法: pug <filename.dog> [-o output]")
		fmt.Println("🔧 オプション:")
		fmt.Println("  --emit-asm    アセンブリコードを表示")
		fmt.Println("  --emit-ast    AST構造を表示")
		fmt.Println("  -O0,-O1,-O2   最適化レベル")
		os.Exit(1)
	}

	filename := os.Args[1]
	fmt.Printf("📄 ファイル '%s' をコンパイル中...\n", filename)

	// ファイルを読み込み
	// #nosec G304 - コンパイラツールとしてファイル読み込みは必要な機能
	input, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("❌ ファイル読み込みエラー: %v\n", err)
		os.Exit(1)
	}

	// 字句解析
	lexer := phase1.New(string(input))

	// 構文解析
	parser := phase1.NewParser(lexer)
	program := parser.ParseProgram()

	// パースエラーをチェック
	errors := parser.Errors()
	if len(errors) > 0 {
		fmt.Println("❌ パースエラー:")
		for _, err := range errors {
			fmt.Printf("  %s\n", err)
		}
		os.Exit(1)
	}

	// コード生成
	codegen := phase2.NewCodeGenerator()
	asmCode, err := codegen.Generate(program)
	if err != nil {
		fmt.Printf("❌ コード生成エラー: %v\n", err)
		os.Exit(1)
	}

	// アセンブリコードを表示
	fmt.Println("✅ アセンブリコード生成成功:")
	fmt.Println(asmCode)
}
