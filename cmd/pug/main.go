package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("🐶 pug コンパイラ - Phase 2 コンパイラ")
	fmt.Println("段階的に学ぶコンパイラ実装プロジェクト")

	if len(os.Args) > 1 {
		filename := os.Args[1]
		fmt.Printf("📄 ファイル '%s' をコンパイル中...\n", filename)
		fmt.Println("⚠️  まだ実装されていません - Phase 2.0 コード生成器から開始してください")
	} else {
		fmt.Println("📝 使用方法: pug <filename.dog> [-o output]")
		fmt.Println("🔧 オプション:")
		fmt.Println("  --emit-asm    アセンブリコードを表示")
		fmt.Println("  --emit-ast    AST構造を表示")
		fmt.Println("  -O0,-O1,-O2   最適化レベル")
	}
}
