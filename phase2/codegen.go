package phase2

import (
	"fmt"
	"strings"

	"github.com/nyasuto/pug/phase1"
)

// CodeGenerator はASTからx86_64アセンブリコードを生成する
type CodeGenerator struct {
	output       strings.Builder
	labelCounter int
	stackOffset  int
	variables    map[string]int // 変数名とスタックオフセットのマッピング
	loopContext  *LoopContext   // 現在のループコンテキスト
}

// NewCodeGenerator は新しいコード生成器を作成する
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		labelCounter: 0,
		stackOffset:  0,
		variables:    make(map[string]int),
	}
}

// Generate はプログラム全体のアセンブリコードを生成する
func (cg *CodeGenerator) Generate(program *phase1.Program) (string, error) {
	// アセンブリのプリアンブル
	cg.emitHeader()

	// 各文を処理
	for _, stmt := range program.Statements {
		if err := cg.generateStatement(stmt); err != nil {
			return "", err
		}
	}

	// アセンブリの後処理
	cg.emitFooter()

	return cg.output.String(), nil
}

// emitHeader はアセンブリファイルのヘッダーを出力する
func (cg *CodeGenerator) emitHeader() {
	cg.emit("# pug compiler generated assembly")
	cg.emit(".section __DATA,__data")
	cg.emit("")
	cg.emit(".section __TEXT,__text,regular,pure_instructions")
	cg.emit(".globl _main")
	cg.emit("")
	cg.emit("_main:")
	cg.emit("    pushq %rbp")      // フレームポインタを保存
	cg.emit("    movq %rsp, %rbp") // 新しいフレームポインタを設定
	cg.emit("    subq $256, %rsp") // ローカル変数用のスタック領域を確保
}

// emitFooter はアセンブリファイルのフッターを出力する
func (cg *CodeGenerator) emitFooter() {
	cg.emit("    movq $0, %rax")   // 戻り値を0に設定
	cg.emit("    movq %rbp, %rsp") // スタックポインタを復元
	cg.emit("    popq %rbp")       // フレームポインタを復元
	cg.emit("    ret")             // 関数から戻る
}

// emit は1行のアセンブリコードを出力する
func (cg *CodeGenerator) emit(code string) {
	cg.output.WriteString(code + "\n")
}

// emitf はフォーマット済みのアセンブリコードを出力する
func (cg *CodeGenerator) emitf(format string, args ...interface{}) {
	cg.emit(fmt.Sprintf(format, args...))
}

// generateLabel は新しいラベルを生成する
func (cg *CodeGenerator) generateLabel(prefix string) string {
	label := fmt.Sprintf(".L%s%d", prefix, cg.labelCounter)
	cg.labelCounter++
	return label
}

// generateStatement は文のアセンブリコードを生成する
func (cg *CodeGenerator) generateStatement(stmt phase1.Statement) error {
	switch node := stmt.(type) {
	case *phase1.LetStatement:
		return cg.generateLetStatement(node)
	case *phase1.ReturnStatement:
		return cg.generateReturnStatement(node)
	case *phase1.ExpressionStatement:
		return cg.generateExpressionStatement(node)
	case *phase1.WhileStatement:
		return cg.generateWhileStatement(node)
	case *phase1.ForStatement:
		return cg.generateForStatement(node)
	case *phase1.BreakStatement:
		return cg.generateBreakStatement(node)
	case *phase1.ContinueStatement:
		return cg.generateContinueStatement(node)
	case *phase1.BlockStatement:
		return cg.generateBlockStatement(node)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

// generateLetStatement はlet文のアセンブリコードを生成する
func (cg *CodeGenerator) generateLetStatement(stmt *phase1.LetStatement) error {
	// 値を計算してRAXに格納
	if err := cg.generateExpression(stmt.Value); err != nil {
		return err
	}

	// 変数をスタックに保存
	cg.stackOffset += 8
	cg.variables[stmt.Name.Value] = cg.stackOffset
	cg.emitf("    movq %%rax, -%d(%%rbp)", cg.stackOffset)
	cg.emitf("    # let %s = ...", stmt.Name.Value)

	return nil
}

// generateReturnStatement はreturn文のアセンブリコードを生成する
func (cg *CodeGenerator) generateReturnStatement(stmt *phase1.ReturnStatement) error {
	if stmt.ReturnValue != nil {
		// 戻り値を計算してRAXに格納
		if err := cg.generateExpression(stmt.ReturnValue); err != nil {
			return err
		}
	} else {
		// 戻り値がない場合は0を返す
		cg.emit("    movq $0, %rax")
	}

	// 関数から戻る
	cg.emit("    movq %rbp, %rsp")
	cg.emit("    popq %rbp")
	cg.emit("    ret")

	return nil
}

// generateExpressionStatement は式文のアセンブリコードを生成する
func (cg *CodeGenerator) generateExpressionStatement(stmt *phase1.ExpressionStatement) error {
	return cg.generateExpression(stmt.Expression)
}

// generateExpression は式のアセンブリコードを生成する
// 結果はRAXレジスタに格納される
func (cg *CodeGenerator) generateExpression(expr phase1.Expression) error {
	switch node := expr.(type) {
	case *phase1.IntegerLiteral:
		return cg.generateIntegerLiteral(node)
	case *phase1.Boolean:
		return cg.generateBoolean(node)
	case *phase1.Identifier:
		return cg.generateIdentifier(node)
	case *phase1.InfixExpression:
		return cg.generateInfixExpression(node)
	case *phase1.PrefixExpression:
		return cg.generatePrefixExpression(node)
	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// generateIntegerLiteral は整数リテラルのアセンブリコードを生成する
func (cg *CodeGenerator) generateIntegerLiteral(node *phase1.IntegerLiteral) error {
	cg.emitf("    movq $%d, %%rax", node.Value)
	return nil
}

// generateBoolean はブール値リテラルのアセンブリコードを生成する
func (cg *CodeGenerator) generateBoolean(node *phase1.Boolean) error {
	if node.Value {
		cg.emit("    movq $1, %rax") // true = 1
	} else {
		cg.emit("    movq $0, %rax") // false = 0
	}
	return nil
}

// generateIdentifier は識別子（変数）のアセンブリコードを生成する
func (cg *CodeGenerator) generateIdentifier(node *phase1.Identifier) error {
	offset, exists := cg.variables[node.Value]
	if !exists {
		return fmt.Errorf("undefined variable: %s", node.Value)
	}

	cg.emitf("    movq -%d(%%rbp), %%rax", offset)
	cg.emitf("    # load variable %s", node.Value)

	return nil
}

// generateInfixExpression は中置式のアセンブリコードを生成する
func (cg *CodeGenerator) generateInfixExpression(node *phase1.InfixExpression) error {
	// 左辺を評価してRAXに格納
	if err := cg.generateExpression(node.Left); err != nil {
		return err
	}

	// RAXをスタックに退避
	cg.emit("    pushq %rax")

	// 右辺を評価してRAXに格納
	if err := cg.generateExpression(node.Right); err != nil {
		return err
	}

	// 右辺の値をRBXに移動
	cg.emit("    movq %rax, %rbx")

	// 左辺の値をスタックから復元
	cg.emit("    popq %rax")

	// 演算子に応じてコードを生成
	switch node.Operator {
	case "+":
		cg.emit("    addq %rbx, %rax")
	case "-":
		cg.emit("    subq %rbx, %rax")
	case "*":
		cg.emit("    imulq %rbx, %rax")
	case "/":
		// 符号拡張してRDX:RAXに展開
		cg.emit("    cqto")
		// RBXで除算（商はRAX、余りはRDX）
		cg.emit("    idivq %rbx")
	case "%":
		// 符号拡張してRDX:RAXに展開
		cg.emit("    cqto")
		// RBXで除算（商はRAX、余りはRDX）
		cg.emit("    idivq %rbx")
		// 余りをRAXに移動
		cg.emit("    movq %rdx, %rax")
	case "==":
		return cg.generateComparison(node.Operator)
	case "!=":
		return cg.generateComparison(node.Operator)
	case "<":
		return cg.generateComparison(node.Operator)
	case ">":
		return cg.generateComparison(node.Operator)
	default:
		return fmt.Errorf("unsupported infix operator: %s", node.Operator)
	}

	return nil
}

// generateComparison は比較演算のアセンブリコードを生成する
func (cg *CodeGenerator) generateComparison(operator string) error {
	// cmpq命令で比較
	cg.emit("    cmpq %rbx, %rax")

	trueLabel := cg.generateLabel("true")
	endLabel := cg.generateLabel("end")

	// 条件に応じてジャンプ
	switch operator {
	case "==":
		cg.emitf("    je %s", trueLabel)
	case "!=":
		cg.emitf("    jne %s", trueLabel)
	case "<":
		cg.emitf("    jl %s", trueLabel)
	case ">":
		cg.emitf("    jg %s", trueLabel)
	}

	// 偽の場合：0をRAXに設定
	cg.emit("    movq $0, %rax")
	cg.emitf("    jmp %s", endLabel)

	// 真の場合：1をRAXに設定
	cg.emitf("%s:", trueLabel)
	cg.emit("    movq $1, %rax")

	cg.emitf("%s:", endLabel)

	return nil
}

// generatePrefixExpression は前置式のアセンブリコードを生成する
func (cg *CodeGenerator) generatePrefixExpression(node *phase1.PrefixExpression) error {
	// 右辺を評価
	if err := cg.generateExpression(node.Right); err != nil {
		return err
	}

	switch node.Operator {
	case "-":
		cg.emit("    negq %rax")
	case "!":
		// ブール値の反転
		cg.emit("    testq %rax, %rax")
		trueLabel := cg.generateLabel("true")
		endLabel := cg.generateLabel("end")

		cg.emitf("    jz %s", trueLabel) // 0の場合は真（!false = true）
		cg.emit("    movq $0, %rax")     // 非0の場合は偽（!true = false）
		cg.emitf("    jmp %s", endLabel)

		cg.emitf("%s:", trueLabel)
		cg.emit("    movq $1, %rax")

		cg.emitf("%s:", endLabel)
	default:
		return fmt.Errorf("unsupported prefix operator: %s", node.Operator)
	}

	return nil
}

// generateWhileStatement はwhile文のアセンブリコードを生成する
func (cg *CodeGenerator) generateWhileStatement(stmt *phase1.WhileStatement) error {
	startLabel := cg.generateLabel("while_start")
	endLabel := cg.generateLabel("while_end")

	// ループコンテキストを設定
	oldContext := cg.loopContext
	cg.loopContext = NewLoopContext(endLabel, startLabel, oldContext)

	// ループ開始ラベル
	cg.emitf("%s:", startLabel)

	// 条件式を評価
	if err := cg.generateExpression(stmt.Condition); err != nil {
		return err
	}

	// 条件が偽の場合はループを抜ける
	cg.emit("    testq %rax, %rax")
	cg.emitf("    jz %s", endLabel)

	// ループ本体を実行
	if err := cg.generateBlockStatement(stmt.Body); err != nil {
		return err
	}

	// ループ開始に戻る
	cg.emitf("    jmp %s", startLabel)

	// ループ終了ラベル
	cg.emitf("%s:", endLabel)

	// ループコンテキストを復元
	cg.loopContext = oldContext

	return nil
}

// generateForStatement はfor文のアセンブリコードを生成する
func (cg *CodeGenerator) generateForStatement(stmt *phase1.ForStatement) error {
	startLabel := cg.generateLabel("for_start")
	continueLabel := cg.generateLabel("for_continue")
	endLabel := cg.generateLabel("for_end")

	// ループコンテキストを設定
	oldContext := cg.loopContext
	cg.loopContext = NewLoopContext(endLabel, continueLabel, oldContext)

	// 初期化文を実行
	if stmt.Initializer != nil {
		if err := cg.generateStatement(stmt.Initializer); err != nil {
			return err
		}
	}

	// ループ開始ラベル
	cg.emitf("%s:", startLabel)

	// 条件式がある場合は評価
	if stmt.Condition != nil {
		if err := cg.generateExpression(stmt.Condition); err != nil {
			return err
		}
		// 条件が偽の場合はループを抜ける
		cg.emit("    testq %rax, %rax")
		cg.emitf("    jz %s", endLabel)
	}

	// ループ本体を実行
	if err := cg.generateBlockStatement(stmt.Body); err != nil {
		return err
	}

	// continueラベル（continue文がここにジャンプ）
	cg.emitf("%s:", continueLabel)

	// 更新式を実行
	if stmt.Update != nil {
		if err := cg.generateExpression(stmt.Update); err != nil {
			return err
		}
	}

	// ループ開始に戻る
	cg.emitf("    jmp %s", startLabel)

	// ループ終了ラベル
	cg.emitf("%s:", endLabel)

	// ループコンテキストを復元
	cg.loopContext = oldContext

	return nil
}

// generateBreakStatement はbreak文のアセンブリコードを生成する
func (cg *CodeGenerator) generateBreakStatement(_ *phase1.BreakStatement) error {
	if cg.loopContext == nil {
		return fmt.Errorf("break statement outside of loop")
	}

	// 現在のループの終了ラベルにジャンプ
	cg.emitf("    jmp %s", cg.loopContext.BreakLabel)
	return nil
}

// generateContinueStatement はcontinue文のアセンブリコードを生成する
func (cg *CodeGenerator) generateContinueStatement(_ *phase1.ContinueStatement) error {
	if cg.loopContext == nil {
		return fmt.Errorf("continue statement outside of loop")
	}

	// 現在のループの継続ラベルにジャンプ
	cg.emitf("    jmp %s", cg.loopContext.ContinueLabel)
	return nil
}

// generateBlockStatement はブロック文のアセンブリコードを生成する
func (cg *CodeGenerator) generateBlockStatement(stmt *phase1.BlockStatement) error {
	for _, s := range stmt.Statements {
		if err := cg.generateStatement(s); err != nil {
			return err
		}
	}
	return nil
}

// AssembleAndLink は生成されたアセンブリコードをアセンブルしてリンクする
func (cg *CodeGenerator) AssembleAndLink(asmCode, outputFile string) error {
	// 実装は後で追加予定
	// アセンブラ（as）とリンカ（ld）を使用して実行可能ファイルを生成
	return fmt.Errorf("assembleAndLink not implemented yet")
}
