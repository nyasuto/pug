package phase1

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const PROMPT = "pug> "

// Start はREPLを開始する
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := NewEnvironment()
	history := []string{}

	// ウェルカムメッセージとヘルプ
	fmt.Fprintln(out, "🔄 Pug Interactive REPL")
	fmt.Fprintln(out, "Type ':help' for help, ':history' for command history, ':exit' to quit")
	fmt.Fprintln(out)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := strings.TrimSpace(scanner.Text())

		// 空行をスキップ
		if line == "" {
			continue
		}

		// 特殊コマンドの処理
		if strings.HasPrefix(line, ":") {
			if handleSpecialCommand(line, out, history) {
				return // :exit コマンドの場合
			}
			continue
		}

		// 履歴に追加
		history = append(history, line)

		// 通常の実行
		l := New(line)
		p := NewParser(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := Eval(program, env)
		if evaluated != nil {
			// エラーの場合は特別な表示
			if evaluated.Type() == ERROR_OBJ {
				fmt.Fprintf(out, "❌ エラー: %s\n", evaluated.Inspect())
			} else if evaluated != NULL_OBJ_INSTANCE {
				fmt.Fprintf(out, "=> %s\n", evaluated.Inspect())
			}
		}
	}
}

// handleSpecialCommand は特殊コマンドを処理する
func handleSpecialCommand(command string, out io.Writer, history []string) bool {
	switch command {
	case ":exit", ":quit", ":q":
		fmt.Fprintln(out, "👋 Goodbye!")
		return true
	case ":help", ":h":
		printHelp(out)
	case ":history":
		printHistory(out, history)
	case ":env":
		fmt.Fprintln(out, "📊 現在の環境変数:")
		fmt.Fprintln(out, "（実装予定: 環境の内容表示）")
	case ":clear":
		// クリア（簡易実装）
		for i := 0; i < 50; i++ {
			fmt.Fprintln(out)
		}
		fmt.Fprintln(out, "🔄 Pug Interactive REPL")
	default:
		fmt.Fprintf(out, "❓ 不明なコマンド: %s\n", command)
		fmt.Fprintln(out, "Type ':help' for available commands")
	}
	return false
}

// printHelp はヘルプメッセージを表示する
func printHelp(out io.Writer) {
	fmt.Fprintln(out, "📖 Pug REPL ヘルプ:")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "🔧 特殊コマンド:")
	fmt.Fprintln(out, "  :help, :h       - このヘルプを表示")
	fmt.Fprintln(out, "  :history        - コマンド履歴を表示")
	fmt.Fprintln(out, "  :env            - 現在の環境変数を表示")
	fmt.Fprintln(out, "  :clear          - 画面をクリア")
	fmt.Fprintln(out, "  :exit, :quit, :q - REPLを終了")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "📝 使用例:")
	fmt.Fprintln(out, "  5 + 3                    # 算術演算")
	fmt.Fprintln(out, "  let x = 42               # 変数代入")
	fmt.Fprintln(out, "  puts(\"Hello World\")      # 出力")
	fmt.Fprintln(out, "  let f = fn(x) { x * 2 }  # 関数定義")
	fmt.Fprintln(out, "  f(10)                    # 関数呼び出し")
	fmt.Fprintln(out, "")
}

// printHistory はコマンド履歴を表示する
func printHistory(out io.Writer, history []string) {
	if len(history) == 0 {
		fmt.Fprintln(out, "📜 履歴: (空)")
		return
	}

	fmt.Fprintln(out, "📜 コマンド履歴:")
	for i, cmd := range history {
		fmt.Fprintf(out, "  %d: %s\n", i+1, cmd)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprintln(out, "🚨 構文解析エラー:")
	for _, msg := range errors {
		fmt.Fprintf(out, "  %s\n", msg)
	}
}
