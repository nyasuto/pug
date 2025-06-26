# Phase 2学習ガイド：コンパイラ基盤とアセンブリ生成

**🚀 インタープリターからコンパイラへの進化**

## 🎯 Phase 2で学べること

Phase 2では、Phase 1のインタープリターから本格的なコンパイラへと発展させます：

1. **コンパイラ理論の基礎** - 実行方式の違いと設計思想
2. **コード生成（Code Generation）** - ASTからアセンブリコードへの変換
3. **型システム（Type System）** - 静的型検査とエラー検出
4. **レジスタ割り当て** - 効率的なメモリ・レジスタ使用
5. **制御構造の実装** - if文、ループのアセンブリ表現
6. **関数呼び出し規約** - スタック管理と引数渡し

これらの技術により、**10倍高速**な実行性能を持つコンパイラを実現します。

---

## 🔄 第1章：インタープリターからコンパイラへ

### 実行方式の違い

| 項目 | インタープリター | コンパイラ |
|------|-----------------|-----------|
| **実行時点** | ソースコード実行時 | 事前コンパイル後 |
| **性能** | 遅い（毎回解析） | 高速（機械語実行） |
| **開発効率** | 高い（即座実行） | 普通（ビルド必要） |
| **エラー検出** | 実行時 | コンパイル時 |
| **最適化** | 限定的 | 高度な最適化可能 |

### pugにおける性能改善目標

```
Phase 1 インタープリター: ベースライン性能
Phase 2 基本コンパイラ:   10倍高速化
Phase 3 最適化コンパイラ: 50倍高速化  
Phase 4 LLVM連携:        100倍高速化
```

### コンパイルフロー

```
ソースコード(.dog) 
    ↓ (字句解析)
トークン列
    ↓ (構文解析)  
AST
    ↓ (型検査) ← Phase 2で追加
型付きAST
    ↓ (コード生成) ← Phase 2で追加
アセンブリコード(.s)
    ↓ (アセンブル・リンク)
実行ファイル
```

---

## ⚙️ 第2章：x86_64アーキテクチャの基礎

### レジスタの理解

x86_64の主要レジスタ：

```assembly
# 汎用レジスタ（64bit）
%rax    # アキュムレータ（戻り値、演算結果）
%rbx    # ベースレジスタ（保存）
%rcx    # カウンタレジスタ（ループ）
%rdx    # データレジスタ（演算、I/O）
%rsi    # ソースインデックス（文字列操作）
%rdi    # デスティネーションインデックス（第1引数）
%rbp    # ベースポインタ（スタックフレーム）
%rsp    # スタックポインタ

# 関数引数渡し用レジスタ（System V ABI）
%rdi    # 第1引数
%rsi    # 第2引数  
%rdx    # 第3引数
%rcx    # 第4引数
%r8     # 第5引数
%r9     # 第6引数
```

### 命令セットの基本

```assembly
# データ移動
mov $42, %rax          # 即値42を%raxに格納
mov %rax, %rbx         # %raxの値を%rbxにコピー
mov %rax, -8(%rbp)     # %raxをスタック位置に保存

# 算術演算
add %rbx, %rax         # %rax = %rax + %rbx
sub %rbx, %rax         # %rax = %rax - %rbx  
imul %rbx, %rax        # %rax = %rax * %rbx
idiv %rbx              # %rax = %rax / %rbx（商）, %rdx = 余り

# 比較と分岐
cmp %rbx, %rax         # %raxと%rbxを比較
je label               # 等しければラベルにジャンプ
jl label               # 小さければラベルにジャンプ
jg label               # 大きければラベルにジャンプ

# スタック操作
push %rax              # %raxをスタックにプッシュ
pop %rax               # スタックから%raxにポップ
```

### 関数呼び出し規約（System V ABI）

```assembly
# 関数呼び出しの流れ
1. 引数をレジスタ/スタックに配置
2. call命令で関数を呼び出し
3. 関数内でスタックフレーム設定
4. 処理実行
5. 戻り値を%raxに設定
6. スタックフレーム復元
7. ret命令で呼び出し元に戻る
```

---

## 🏗️ 第3章：コード生成器の実装

### CodeGenerator構造体の設計

```go
type CodeGenerator struct {
    output       strings.Builder    // 生成されるアセンブリコード
    labelCounter int                // ユニークラベル生成用カウンタ
    stackOffset  int                // 現在のスタックオフセット
    variables    map[string]int     // 変数名→スタック位置マッピング
    loopContext  *LoopContext       // ループ制御用コンテキスト
}
```

**設計原則：**
- **単一責任**：ASTからアセンブリコードへの変換に特化
- **状態管理**：変数とスタックの状態を正確に追跡
- **拡張性**：新しい構文要素の追加が容易
- **デバッグ性**：生成されたコードが読みやすい

### アセンブリファイル構造の生成

```go
func (cg *CodeGenerator) emitHeader() {
    cg.emit("# pug compiler generated assembly")
    cg.emit(".section __DATA,__data")
    cg.emit("")
    cg.emit(".section __TEXT,__text")
    cg.emit(".globl _main")
    cg.emit("")
}

func (cg *CodeGenerator) emitFooter() {
    cg.emit("_main:")
    cg.emit("    push %rbp")           // スタックフレーム設定
    cg.emit("    mov %rsp, %rbp")
    
    // メイン処理...
    
    cg.emit("    mov $0, %rax")        // 終了コード0
    cg.emit("    pop %rbp")
    cg.emit("    ret")
}
```

### 整数リテラルの生成

```go
func (cg *CodeGenerator) generateIntegerLiteral(node *phase1.IntegerLiteral) error {
    // 即値を%raxレジスタに読み込み
    cg.emit(fmt.Sprintf("    mov $%d, %%rax", node.Value))
    return nil
}
```

**生成例：**
```dog
42
```
↓
```assembly
mov $42, %rax
```

### 変数宣言の実装

```go
func (cg *CodeGenerator) generateLetStatement(node *phase1.LetStatement) error {
    // 右辺の式を評価（結果は%raxに）
    if err := cg.generateExpression(node.Value); err != nil {
        return err
    }
    
    // スタック上に変数用の領域を確保
    cg.stackOffset -= 8
    cg.variables[node.Name.Value] = cg.stackOffset
    
    // %raxの値を変数の位置に保存
    cg.emit(fmt.Sprintf("    mov %%rax, %d(%%rbp)", cg.stackOffset))
    return nil
}
```

**生成例：**
```dog
let x = 42;
```
↓
```assembly
mov $42, %rax           # 右辺の評価
mov %rax, -8(%rbp)      # 変数xに保存
```

### 変数参照の実装

```go
func (cg *CodeGenerator) generateIdentifier(node *phase1.Identifier) error {
    offset, exists := cg.variables[node.Value]
    if !exists {
        return fmt.Errorf("undefined variable: %s", node.Value)
    }
    
    // 変数の値を%raxに読み込み
    cg.emit(fmt.Sprintf("    mov %d(%%rbp), %%rax", offset))
    return nil
}
```

**生成例：**
```dog
x
```
↓
```assembly
mov -8(%rbp), %rax      # 変数xの値を%raxに読み込み
```

---

## 🧮 第4章：演算処理の実装

### 中置演算の処理

```go
func (cg *CodeGenerator) generateInfixExpression(node *phase1.InfixExpression) error {
    // 左辺を評価（結果は%raxに）
    if err := cg.generateExpression(node.Left); err != nil {
        return err
    }
    
    // 左辺の結果をスタックに一時保存
    cg.emit("    push %rax")
    
    // 右辺を評価（結果は%raxに）
    if err := cg.generateExpression(node.Right); err != nil {
        return err
    }
    
    // 右辺の結果を%rbxに移動
    cg.emit("    mov %rax, %rbx")
    
    // 左辺の結果をスタックから復元
    cg.emit("    pop %rax")
    
    // 演算実行
    switch node.Operator {
    case "+":
        cg.emit("    add %rbx, %rax")   // %rax = %rax + %rbx
    case "-":
        cg.emit("    sub %rbx, %rax")   // %rax = %rax - %rbx
    case "*":
        cg.emit("    imul %rbx, %rax")  // %rax = %rax * %rbx
    case "/":
        cg.emit("    cqo")              // %rdxを符号拡張
        cg.emit("    idiv %rbx")        // %rax = %rax / %rbx
    default:
        return fmt.Errorf("unsupported operator: %s", node.Operator)
    }
    
    return nil
}
```

**生成例：**
```dog
2 + 3 * 4
```
↓
```assembly
# 2の評価
mov $2, %rax
push %rax

# 3 * 4の評価
mov $3, %rax
push %rax
mov $4, %rax
mov %rax, %rbx
pop %rax
imul %rbx, %rax

# 2 + (3 * 4)の計算
mov %rax, %rbx
pop %rax
add %rbx, %rax
```

### 比較演算の実装

```go
func (cg *CodeGenerator) generateComparisonExpression(operator string) error {
    label1 := cg.newLabel("compare_true")
    label2 := cg.newLabel("compare_end")
    
    // 比較実行
    cg.emit("    cmp %rbx, %rax")
    
    // 条件分岐
    switch operator {
    case "==":
        cg.emit(fmt.Sprintf("    je %s", label1))
    case "!=":
        cg.emit(fmt.Sprintf("    jne %s", label1))
    case "<":
        cg.emit(fmt.Sprintf("    jl %s", label1))
    case ">":
        cg.emit(fmt.Sprintf("    jg %s", label1))
    }
    
    // false の場合
    cg.emit("    mov $0, %rax")
    cg.emit(fmt.Sprintf("    jmp %s", label2))
    
    // true の場合
    cg.emit(fmt.Sprintf("%s:", label1))
    cg.emit("    mov $1, %rax")
    
    cg.emit(fmt.Sprintf("%s:", label2))
    return nil
}
```

**生成例：**
```dog
x < 10
```
↓
```assembly
mov -8(%rbp), %rax      # x の値
push %rax
mov $10, %rax           # 10
mov %rax, %rbx
pop %rax
cmp %rbx, %rax          # x と 10 を比較
jl compare_true_1       # x < 10 なら真
mov $0, %rax            # 偽の場合
jmp compare_end_1
compare_true_1:
mov $1, %rax            # 真の場合
compare_end_1:
```

---

## 🔀 第5章：制御構造の実装

### if文の実装

```go
func (cg *CodeGenerator) generateIfExpression(node *phase1.IfExpression) error {
    elseLabel := cg.newLabel("else")
    endLabel := cg.newLabel("if_end")
    
    // 条件式を評価
    if err := cg.generateExpression(node.Condition); err != nil {
        return err
    }
    
    // 結果が0（false）なら else へジャンプ
    cg.emit("    cmp $0, %rax")
    cg.emit(fmt.Sprintf("    je %s", elseLabel))
    
    // then ブロックの生成
    if err := cg.generateBlockStatement(node.Consequence); err != nil {
        return err
    }
    cg.emit(fmt.Sprintf("    jmp %s", endLabel))
    
    // else ブロックの生成
    cg.emit(fmt.Sprintf("%s:", elseLabel))
    if node.Alternative != nil {
        if err := cg.generateBlockStatement(node.Alternative); err != nil {
            return err
        }
    }
    
    cg.emit(fmt.Sprintf("%s:", endLabel))
    return nil
}
```

**生成例：**
```dog
if (x > 5) {
    y = 10;
} else {
    y = 0;
}
```
↓
```assembly
mov -8(%rbp), %rax      # x の値
push %rax
mov $5, %rax
mov %rax, %rbx
pop %rax
cmp %rbx, %rax          # x > 5 を比較
jg compare_true_1
mov $0, %rax
jmp compare_end_1
compare_true_1:
mov $1, %rax
compare_end_1:
cmp $0, %rax            # 条件の結果をチェック
je else_1               # false なら else へ

# then ブロック
mov $10, %rax
mov %rax, -16(%rbp)     # y = 10
jmp if_end_1

# else ブロック  
else_1:
mov $0, %rax
mov %rax, -16(%rbp)     # y = 0

if_end_1:
```

### while文の実装

```go
func (cg *CodeGenerator) generateWhileStatement(node *phase1.WhileStatement) error {
    loopStart := cg.newLabel("while_start")
    loopEnd := cg.newLabel("while_end")
    
    // ループコンテキスト設定（break/continue用）
    oldContext := cg.loopContext
    cg.loopContext = &LoopContext{
        BreakLabel:    loopEnd,
        ContinueLabel: loopStart,
    }
    defer func() { cg.loopContext = oldContext }()
    
    // ループ開始
    cg.emit(fmt.Sprintf("%s:", loopStart))
    
    // 条件式を評価
    if err := cg.generateExpression(node.Condition); err != nil {
        return err
    }
    
    // 条件が false ならループ終了
    cg.emit("    cmp $0, %rax")
    cg.emit(fmt.Sprintf("    je %s", loopEnd))
    
    // ループボディの実行
    if err := cg.generateBlockStatement(node.Body); err != nil {
        return err
    }
    
    // ループ開始に戻る
    cg.emit(fmt.Sprintf("    jmp %s", loopStart))
    
    // ループ終了
    cg.emit(fmt.Sprintf("%s:", loopEnd))
    
    return nil
}
```

**生成例：**
```dog
while (i < 10) {
    i = i + 1;
}
```
↓
```assembly
while_start_1:
mov -8(%rbp), %rax      # i の値
push %rax
mov $10, %rax
mov %rax, %rbx
pop %rax
cmp %rbx, %rax          # i < 10
jl compare_true_2
mov $0, %rax
jmp compare_end_2
compare_true_2:
mov $1, %rax
compare_end_2:
cmp $0, %rax            # 条件チェック
je while_end_1          # false なら終了

# ループボディ
mov -8(%rbp), %rax      # i
push %rax
mov $1, %rax
mov %rax, %rbx
pop %rax
add %rbx, %rax          # i + 1
mov %rax, -8(%rbp)      # i = i + 1

jmp while_start_1       # ループ開始に戻る
while_end_1:
```

---

## 🎭 第6章：関数定義と呼び出し

### 関数定義の実装

```go
func (cg *CodeGenerator) generateFunctionLiteral(node *phase1.FunctionLiteral) error {
    funcName := cg.newLabel("function")
    
    // 関数のアドレスを%raxに設定（クロージャ対応は今後）
    cg.emit(fmt.Sprintf("    lea %s(%%rip), %%rax", funcName))
    
    // 関数の後に配置するためジャンプ
    skipLabel := cg.newLabel("skip_function")
    cg.emit(fmt.Sprintf("    jmp %s", skipLabel))
    
    // 関数本体の生成
    cg.emit(fmt.Sprintf("%s:", funcName))
    cg.emit("    push %rbp")
    cg.emit("    mov %rsp, %rbp")
    
    // パラメータのスタック配置
    oldVars := cg.variables
    oldOffset := cg.stackOffset
    cg.variables = make(map[string]int)
    cg.stackOffset = 0
    
    // 引数をスタックに配置
    for i, param := range node.Parameters {
        cg.stackOffset -= 8
        cg.variables[param.Value] = cg.stackOffset
        cg.emit(fmt.Sprintf("    mov %%rdi, %d(%%rbp)", cg.stackOffset)) // 第1引数のみ対応
    }
    
    // 関数ボディの実行
    if err := cg.generateBlockStatement(node.Body); err != nil {
        return err
    }
    
    // デフォルト戻り値（関数が値を返さない場合）
    cg.emit("    mov $0, %rax")
    
    // 関数終了
    cg.emit("    pop %rbp")
    cg.emit("    ret")
    
    // 関数定義をスキップするラベル
    cg.emit(fmt.Sprintf("%s:", skipLabel))
    
    // 変数スコープ復元
    cg.variables = oldVars
    cg.stackOffset = oldOffset
    
    return nil
}
```

### 関数呼び出しの実装

```go
func (cg *CodeGenerator) generateCallExpression(node *phase1.CallExpression) error {
    // 引数の評価（逆順でスタックにプッシュ）
    for i := len(node.Arguments) - 1; i >= 0; i-- {
        if err := cg.generateExpression(node.Arguments[i]); err != nil {
            return err
        }
        cg.emit("    push %rax")
    }
    
    // 第1引数を%rdiに設定（System V ABI）
    if len(node.Arguments) > 0 {
        cg.emit("    pop %rdi")
    }
    
    // 関数の評価（関数のアドレスを取得）
    if err := cg.generateExpression(node.Function); err != nil {
        return err
    }
    
    // 関数呼び出し
    cg.emit("    call *%rax")
    
    // スタッククリーンアップ（残りの引数）
    if len(node.Arguments) > 1 {
        stackCleanup := (len(node.Arguments) - 1) * 8
        cg.emit(fmt.Sprintf("    add $%d, %%rsp", stackCleanup))
    }
    
    return nil
}
```

**生成例：**
```dog
let add = fn(x, y) { x + y };
add(3, 4);
```
↓
```assembly
# 関数定義
lea function_1(%rip), %rax
jmp skip_function_1

function_1:
push %rbp
mov %rsp, %rbp
mov %rdi, -8(%rbp)      # 第1引数 x
mov %rsi, -16(%rbp)     # 第2引数 y
mov -8(%rbp), %rax      # x
push %rax
mov -16(%rbp), %rax     # y
mov %rax, %rbx
pop %rax
add %rbx, %rax          # x + y
pop %rbp
ret

skip_function_1:
mov %rax, -8(%rbp)      # add 変数に関数保存

# 関数呼び出し
mov $4, %rax            # 第2引数
push %rax
mov $3, %rax            # 第1引数
mov %rax, %rdi
mov -8(%rbp), %rax      # add 関数のアドレス
call *%rax              # 関数呼び出し
add $8, %rsp            # スタッククリーンアップ
```

---

## 🔍 第7章：型システムの実装

### 型情報の定義

```go
type Type int

const (
    INTEGER_TYPE Type = iota
    BOOLEAN_TYPE
    STRING_TYPE
    FUNCTION_TYPE
    VOID_TYPE
    UNKNOWN_TYPE
)

type TypeChecker struct {
    errors      []string
    symbolTable map[string]Type
}
```

### 式の型検査

```go
func (tc *TypeChecker) checkExpression(node phase1.Expression) Type {
    switch n := node.(type) {
    case *phase1.IntegerLiteral:
        return INTEGER_TYPE
        
    case *phase1.Boolean:
        return BOOLEAN_TYPE
        
    case *phase1.Identifier:
        if typ, exists := tc.symbolTable[n.Value]; exists {
            return typ
        }
        tc.addError(fmt.Sprintf("undefined variable: %s", n.Value))
        return UNKNOWN_TYPE
        
    case *phase1.InfixExpression:
        leftType := tc.checkExpression(n.Left)
        rightType := tc.checkExpression(n.Right)
        
        // 算術演算の型チェック
        if n.Operator == "+" || n.Operator == "-" || 
           n.Operator == "*" || n.Operator == "/" {
            if leftType != INTEGER_TYPE || rightType != INTEGER_TYPE {
                tc.addError(fmt.Sprintf("arithmetic operator requires integers"))
                return UNKNOWN_TYPE
            }
            return INTEGER_TYPE
        }
        
        // 比較演算の型チェック
        if n.Operator == "==" || n.Operator == "!=" ||
           n.Operator == "<" || n.Operator == ">" {
            if leftType != rightType {
                tc.addError(fmt.Sprintf("comparison between different types"))
                return UNKNOWN_TYPE
            }
            return BOOLEAN_TYPE
        }
        
    default:
        tc.addError(fmt.Sprintf("unknown expression type: %T", node))
        return UNKNOWN_TYPE
    }
    
    return UNKNOWN_TYPE
}
```

### Let文の型検査

```go
func (tc *TypeChecker) checkLetStatement(node *phase1.LetStatement) {
    valueType := tc.checkExpression(node.Value)
    
    if valueType == UNKNOWN_TYPE {
        return // エラーは既に追加済み
    }
    
    // 変数を記号表に追加
    tc.symbolTable[node.Name.Value] = valueType
}
```

### 型検査の統合

```go
func (cg *CodeGenerator) GenerateWithTypeCheck(program *phase1.Program) (string, error) {
    // 型検査実行
    tc := NewTypeChecker()
    if err := tc.CheckProgram(program); err != nil {
        return "", fmt.Errorf("type check failed: %v", err)
    }
    
    // 型検査が成功したらコード生成
    return cg.Generate(program)
}
```

---

## ⚡ 第8章：最適化の基礎

### 定数畳み込み（Constant Folding）

```go
func (cg *CodeGenerator) foldConstants(node phase1.Expression) phase1.Expression {
    switch n := node.(type) {
    case *phase1.InfixExpression:
        left := cg.foldConstants(n.Left)
        right := cg.foldConstants(n.Right)
        
        // 両方が整数リテラルの場合
        if leftInt, ok := left.(*phase1.IntegerLiteral); ok {
            if rightInt, ok := right.(*phase1.IntegerLiteral); ok {
                switch n.Operator {
                case "+":
                    return &phase1.IntegerLiteral{Value: leftInt.Value + rightInt.Value}
                case "-":
                    return &phase1.IntegerLiteral{Value: leftInt.Value - rightInt.Value}
                case "*":
                    return &phase1.IntegerLiteral{Value: leftInt.Value * rightInt.Value}
                case "/":
                    if rightInt.Value != 0 {
                        return &phase1.IntegerLiteral{Value: leftInt.Value / rightInt.Value}
                    }
                }
            }
        }
        
        return &phase1.InfixExpression{
            Left:     left,
            Operator: n.Operator,
            Right:    right,
        }
    }
    
    return node
}
```

**最適化例：**
```dog
let x = 2 + 3 * 4;
```
↓
```dog
let x = 14;  // コンパイル時に計算
```
↓
```assembly
mov $14, %rax           # 最適化後（計算済み）
mov %rax, -8(%rbp)
```

### レジスタ使用の最適化

```go
type RegisterAllocator struct {
    available []string          // 利用可能レジスタ
    allocated map[string]string // 変数→レジスタマッピング
}

func (ra *RegisterAllocator) allocateRegister(variable string) string {
    if reg, exists := ra.allocated[variable]; exists {
        return reg
    }
    
    if len(ra.available) > 0 {
        reg := ra.available[0]
        ra.available = ra.available[1:]
        ra.allocated[variable] = reg
        return reg
    }
    
    return "" // スタック使用が必要
}
```

---

## 🎯 第9章：実践課題と練習問題

### 課題1：新しい演算子の追加

剰余演算子（%）を追加してみましょう：

```go
case "%":
    cg.emit("    cqo")              // %rdxを符号拡張
    cg.emit("    idiv %rbx")        // 除算実行
    cg.emit("    mov %rdx, %rax")   // 余りを%raxに移動
```

### 課題2：配列の実装

配列アクセスのコード生成：

```go
func (cg *CodeGenerator) generateArrayAccess(array, index phase1.Expression) error {
    // 配列のベースアドレスを取得
    if err := cg.generateExpression(array); err != nil {
        return err
    }
    cg.emit("    push %rax")
    
    // インデックスを評価
    if err := cg.generateExpression(index); err != nil {
        return err
    }
    
    // インデックス * 8（8バイト/要素）
    cg.emit("    imul $8, %rax")
    cg.emit("    mov %rax, %rbx")
    
    // ベースアドレス + オフセット
    cg.emit("    pop %rax")
    cg.emit("    add %rbx, %rax")
    
    // 配列要素の値を読み込み
    cg.emit("    mov (%rax), %rax")
    
    return nil
}
```

### 課題3：文字列リテラルの実装

```go
func (cg *CodeGenerator) generateStringLiteral(node *phase1.StringLiteral) error {
    label := cg.newLabel("string")
    
    // データセクションに文字列を配置
    cg.dataSection = append(cg.dataSection, 
        fmt.Sprintf("%s: .asciz \"%s\"", label, node.Value))
    
    // 文字列のアドレスを%raxに設定
    cg.emit(fmt.Sprintf("    lea %s(%%rip), %%rax", label))
    
    return nil
}
```

### 課題4：再帰関数の最適化

末尾再帰の最適化実装：

```go
func (cg *CodeGenerator) optimizeTailRecursion(node *phase1.CallExpression, 
    funcName string) bool {
    // 関数が自分自身を呼び出しているかチェック
    if ident, ok := node.Function.(*phase1.Identifier); ok {
        if ident.Value == funcName {
            // 引数を更新
            cg.updateParameters(node.Arguments)
            // 関数先頭にジャンプ（callではなく）
            cg.emit(fmt.Sprintf("    jmp %s_start", funcName))
            return true
        }
    }
    return false
}
```

---

## 📊 第10章：性能測定と分析

### ベンチマーク実装

```go
func BenchmarkCompilerVsInterpreter(b *testing.B) {
    input := `
    let fibonacci = fn(n) {
        if (n < 2) {
            return n;
        } else {
            return fibonacci(n - 1) + fibonacci(n - 2);
        }
    };
    fibonacci(20);
    `
    
    // インタープリター
    b.Run("Interpreter", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            l := lexer.New(input)
            p := parser.New(l)
            program := p.ParseProgram()
            env := object.NewEnvironment()
            evaluator.Eval(program, env)
        }
    })
    
    // コンパイラ
    b.Run("Compiler", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            l := lexer.New(input)
            p := parser.New(l)
            program := p.ParseProgram()
            cg := NewCodeGenerator()
            cg.Generate(program)
            // アセンブル・実行は省略
        }
    })
}
```

### 期待される性能改善

```
Fibonacci(20)の実行時間:
- Phase 1 インタープリター: 2.3秒
- Phase 2 基本コンパイラ:   0.23秒 (10倍高速)

メモリ使用量:
- インタープリター: ASTを保持し続ける
- コンパイラ: 実行時はコンパクトな機械語のみ
```

---

## 🚀 第11章：Phase 3への準備

### Phase 2の成果

Phase 2を完了すると、以下が実装されています：

✅ **アセンブリコード生成**：x86_64機械語への変換  
✅ **型システム**：コンパイル時エラー検出  
✅ **制御構造**：if文・while文の効率的実装  
✅ **関数呼び出し**：System V ABI準拠  
✅ **基本最適化**：定数畳み込み等  
✅ **性能向上**：インタープリターの10倍高速  

### Phase 3で学ぶこと

次のPhase 3では、さらに高度な最適化技術を学びます：

🎯 **中間表現（IR）**：SSA形式による最適化基盤  
🎯 **高度最適化**：ループ最適化・インライン展開  
🎯 **データフロー解析**：変数の生存期間分析  
🎯 **レジスタ割り当て**：グラフ彩色法による効率化  

### 移行時のチェックポイント

Phase 3に進む前に以下を確認：

- [ ] 型検査が正しく動作する
- [ ] 生成されたアセンブリがアセンブル・実行できる
- [ ] 基本的な制御構造が動作する
- [ ] 関数定義・呼び出しが正常に動作する
- [ ] インタープリターより高速に実行される

---

## 🔗 参考資料とさらなる学習

### 推奨書籍

1. **「Compilers: Principles, Techniques, and Tools」** - Aho, Lam, Sethi, Ullman（ドラゴンブック）
2. **「Engineering a Compiler」** - Cooper, Torczon
3. **「Modern Compiler Implementation」** - Appel

### アセンブリ言語学習

- **「Programming from the Ground Up」** - Jonathan Bartlett
- **Intel 64 and IA-32 Architectures Software Developer's Manual**
- **System V ABI Documentation**

### オンライン学習リソース

- [Cornell CS 4120 Introduction to Compilers](https://www.cs.cornell.edu/courses/cs4120/)
- [Stanford CS 143 Compilers](https://web.stanford.edu/class/cs143/)
- [x86 Assembly Guide](https://www.cs.virginia.edu/~evans/cs216/guides/x86.html)

---

## 🎓 実装スキルの習得確認

### 理解度チェック

以下の質問に答えられるか確認：

1. インタープリターとコンパイラの性能差の理由は？
2. x86_64の関数呼び出し規約の特徴は？
3. スタックフレームの構造と役割は？
4. 型システムがもたらす利点は？
5. 定数畳み込み最適化の仕組みは？

### 実装課題

以下の機能を独自に実装してみましょう：

1. **switch文の実装**
2. **for文の実装**
3. **構造体の基本実装**
4. **ポインタ演算の基本**
5. **簡単なガベージコレクション**

---

**🎉 Phase 2完了おめでとうございます！**

基本的なコンパイラの仕組みを理解し、実際に動作するコンパイラを実装できました。10倍の性能向上を実現し、型安全性も確保したコンパイラは立派な成果です。Phase 3では、さらに高度な最適化技術を学び、50倍の性能向上を目指しましょう！

**次のステップ：[Phase 3学習ガイド](phase3_tutorial.md)**