package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nyasuto/pug/phase1"
	"github.com/nyasuto/pug/phase2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "使用法: %s <ソースファイル> [出力ファイル]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  pug言語ソースファイルをx86_64アセンブリにコンパイルします\n")
		fmt.Fprintf(os.Stderr, "\n例:\n")
		fmt.Fprintf(os.Stderr, "  %s program.dog\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s program.dog program.s\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// 基本的なパス検証（セキュリティ対策）
	if strings.Contains(inputFile, "..") {
		fmt.Fprintf(os.Stderr, "❌ エラー: 相対パス指定は許可されていません\n")
		os.Exit(1)
	}

	// 出力ファイル名を決定
	var outputFile string
	if len(os.Args) >= 3 {
		outputFile = os.Args[2]
	} else {
		// 拡張子を .s に変更
		ext := filepath.Ext(inputFile)
		outputFile = strings.TrimSuffix(inputFile, ext) + ".s"
	}

	// ソースファイルを読み込み
	// #nosec G304 -- inputFile is validated above to prevent path traversal
	source, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: ファイル読み込み失敗: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📝 ソースファイル: %s\n", inputFile)
	fmt.Printf("🎯 出力ファイル: %s\n", outputFile)

	// 字句解析
	fmt.Println("🔤 字句解析中...")
	lexer := phase1.New(string(source))

	// 構文解析
	fmt.Println("📜 構文解析中...")
	parser := phase1.NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		fmt.Fprintf(os.Stderr, "❌ 構文解析エラー:\n")
		for _, err := range parser.Errors() {
			fmt.Fprintf(os.Stderr, "  %s\n", err)
		}
		os.Exit(1)
	}

	// コード生成
	fmt.Println("⚙️ アセンブリコード生成中...")
	codeGen := phase2.NewCodeGenerator()
	asmCode, err := codeGen.Generate(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ コード生成エラー: %v\n", err)
		os.Exit(1)
	}

	// アセンブリコードをファイルに出力
	fmt.Println("💾 アセンブリファイル出力中...")
	err = os.WriteFile(outputFile, []byte(asmCode), 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ ファイル出力エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ コンパイル完了！\n")
	fmt.Printf("📊 生成されたアセンブリ: %d行\n", strings.Count(asmCode, "\n"))
	fmt.Printf("\n次のステップ:\n")
	fmt.Printf("  アセンブル: as -arch x86_64 -o %s.o %s\n", strings.TrimSuffix(outputFile, ".s"), outputFile)
	fmt.Printf("  リンク:     ld -o %s -lSystem -syslibroot /Library/Developer/CommandLineTools/SDKs/MacOSX.sdk -e _main -arch x86_64 %s.o\n",
		strings.TrimSuffix(outputFile, ".s"), strings.TrimSuffix(outputFile, ".s"))
	fmt.Printf("  実行:       ./%s\n", strings.TrimSuffix(outputFile, ".s"))
}
