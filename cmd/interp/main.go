package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("🐶 pug コンパイラ - Phase 1 インタープリター")
	fmt.Println("段階的に学ぶコンパイラ実装プロジェクト")

	if len(os.Args) > 1 {
		filename := os.Args[1]
		fmt.Printf("📄 ファイル '%s' を処理中...\n", filename)
		fmt.Println("⚠️  まだ実装されていません - Phase 1.0 レクサーから開始してください")
	} else {
		fmt.Println("📝 使用方法: interp <filename.dog>")
		fmt.Println("🔄 または REPL モード: interp --repl")
	}
}
