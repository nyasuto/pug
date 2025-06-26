# 実践課題とサンプルコード集

**🎯 段階的な実装課題でスキルを定着**

このドキュメントでは、pugプロジェクトの各フェーズで取り組める実践課題と、参考となるサンプルコードを提供します。

## 📚 Phase 1 実践課題

### 課題1-1：字句解析器の拡張

**目標**: 新しいトークン種別を追加し、字句解析器を拡張する

**追加するトークン**:
```go
// 比較演算子
token.LE    // <=
token.GE    // >=
token.AND   // &&
token.OR    // ||

// 代入演算子
token.PLUS_ASSIGN  // +=
token.MINUS_ASSIGN // -=
```

**実装例**:
```go
func (l *Lexer) NextToken() Token {
    switch l.ch {
    case '<':
        if l.peekChar() == '=' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: LE, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(LT, l.ch)
        }
    case '&':
        if l.peekChar() == '&' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: AND, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(ILLEGAL, l.ch)
        }
    }
}
```

### 課題1-2：文字列リテラルの実装

**目標**: 文字列リテラルをサポートし、エスケープシーケンスを処理する

**サンプルコード**:
```dog
let message = "Hello, world!";
let path = "C:\\Users\\name\\file.txt";
let multiline = "Line 1\nLine 2\nLine 3";
```

**実装例**:
```go
func (l *Lexer) readString() string {
    position := l.position + 1
    for {
        l.readChar()
        if l.ch == '"' || l.ch == 0 {
            break
        }
        // エスケープシーケンス処理
        if l.ch == '\\' {
            l.readChar() // エスケープ文字をスキップ
        }
    }
    return l.input[position:l.position]
}

func (e *Evaluator) evalStringLiteral(node *ast.StringLiteral) object.Object {
    // エスケープシーケンスを解釈
    value := strings.ReplaceAll(node.Value, "\\n", "\n")
    value = strings.ReplaceAll(value, "\\t", "\t")
    value = strings.ReplaceAll(value, "\\\\", "\\")
    return &object.String{Value: value}
}
```

### 課題1-3：配列リテラルの実装

**目標**: 配列リテラルと配列アクセスを実装する

**サンプルコード**:
```dog
let numbers = [1, 2, 3, 4, 5];
let first = numbers[0];
let length = len(numbers);
```

**実装例**:
```go
// AST定義
type ArrayLiteral struct {
    Token    token.Token
    Elements []Expression
}

type IndexExpression struct {
    Token token.Token
    Left  Expression
    Index Expression
}

// パーサー実装
func (p *Parser) parseArrayLiteral() ast.Expression {
    array := &ast.ArrayLiteral{Token: p.curToken}
    array.Elements = p.parseExpressionList(token.RBRACKET)
    return array
}

// 評価器実装
func (e *Evaluator) evalArrayLiteral(node *ast.ArrayLiteral, env *object.Environment) object.Object {
    elements := []object.Object{}
    for _, elem := range node.Elements {
        evaluated := e.Eval(elem, env)
        if isError(evaluated) {
            return evaluated
        }
        elements = append(elements, evaluated)
    }
    return &object.Array{Elements: elements}
}
```

### 課題1-4：while文の実装

**目標**: while文を追加してループ処理を可能にする

**サンプルコード**:
```dog
let i = 0;
while (i < 10) {
    print(i);
    i = i + 1;
}
```

**実装例**:
```go
// AST定義
type WhileStatement struct {
    Token     token.Token
    Condition Expression
    Body      *BlockStatement
}

// パーサー実装
func (p *Parser) parseWhileStatement() ast.Statement {
    stmt := &ast.WhileStatement{Token: p.curToken}
    
    if !p.expectPeek(token.LPAREN) {
        return nil
    }
    
    p.nextToken()
    stmt.Condition = p.parseExpression(LOWEST)
    
    if !p.expectPeek(token.RPAREN) {
        return nil
    }
    
    if !p.expectPeek(token.LBRACE) {
        return nil
    }
    
    stmt.Body = p.parseBlockStatement()
    return stmt
}

// 評価器実装
func (e *Evaluator) evalWhileStatement(node *ast.WhileStatement, env *object.Environment) object.Object {
    for {
        condition := e.Eval(node.Condition, env)
        if isError(condition) {
            return condition
        }
        
        if !isTruthy(condition) {
            break
        }
        
        result := e.Eval(node.Body, env)
        if result != nil {
            rt := result.Type()
            if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
                return result
            }
        }
    }
    return NULL
}
```

## 🚀 Phase 2 実践課題

### 課題2-1：複雑な演算子の実装

**目標**: 複合代入演算子とインクリメント/デクリメント演算子を実装

**サンプルコード**:
```dog
let x = 10;
x += 5;    // x = x + 5
x++;       // x = x + 1
++x;       // x = x + 1
```

**アセンブリ生成例**:
```assembly
# x += 5 の生成
mov -8(%rbp), %rax    # x の値をロード
add $5, %rax          # 5を加算
mov %rax, -8(%rbp)    # 結果を x に保存

# x++ の生成
mov -8(%rbp), %rax    # x の値をロード
inc %rax              # インクリメント
mov %rax, -8(%rbp)    # 結果を x に保存
```

### 課題2-2：配列のコード生成

**目標**: 配列アクセスのアセンブリコードを生成する

**実装例**:
```go
func (cg *CodeGenerator) generateIndexExpression(node *ast.IndexExpression) error {
    // 配列のベースアドレスを取得
    if err := cg.generateExpression(node.Left); err != nil {
        return err
    }
    cg.emit("    push %rax")
    
    // インデックスを評価
    if err := cg.generateExpression(node.Index); err != nil {
        return err
    }
    
    // インデックス * 8（64bit要素サイズ）
    cg.emit("    imul $8, %rax")
    cg.emit("    mov %rax, %rbx")
    
    // ベースアドレス + インデックス
    cg.emit("    pop %rax")
    cg.emit("    add %rbx, %rax")
    
    // 配列要素の値をロード
    cg.emit("    mov (%rax), %rax")
    
    return nil
}
```

### 課題2-3：構造体の基本実装

**目標**: 構造体定義とフィールドアクセスを実装する

**サンプルコード**:
```dog
struct Point {
    x: int,
    y: int,
}

let p = Point{x: 10, y: 20};
let x_val = p.x;
```

**実装例**:
```go
// AST定義
type StructLiteral struct {
    Token  token.Token
    Fields map[string]Expression
}

type FieldAccess struct {
    Token token.Token
    Left  Expression
    Field *Identifier
}

// 型システム
type StructType struct {
    Fields map[string]Type
}

func (tc *TypeChecker) checkStructLiteral(node *ast.StructLiteral) Type {
    fieldTypes := make(map[string]Type)
    for name, expr := range node.Fields {
        fieldTypes[name] = tc.checkExpression(expr)
    }
    return &StructType{Fields: fieldTypes}
}
```

### 課題2-4：関数の高度な呼び出し規約

**目標**: 複数引数を持つ関数の正しい呼び出し規約を実装

**System V ABI準拠の実装**:
```go
func (cg *CodeGenerator) generateFunctionCall(node *ast.CallExpression) error {
    args := node.Arguments
    
    // レジスタ引数（最初の6個）
    regArgs := []string{"%rdi", "%rsi", "%rdx", "%rcx", "%r8", "%r9"}
    
    // スタック引数（7個目以降、逆順でプッシュ）
    for i := len(args) - 1; i >= 6; i-- {
        if err := cg.generateExpression(args[i]); err != nil {
            return err
        }
        cg.emit("    push %rax")
    }
    
    // レジスタ引数（正順で設定）
    for i := 0; i < len(args) && i < 6; i++ {
        if err := cg.generateExpression(args[i]); err != nil {
            return err
        }
        cg.emit(fmt.Sprintf("    mov %%rax, %s", regArgs[i]))
    }
    
    // 関数呼び出し
    if err := cg.generateExpression(node.Function); err != nil {
        return err
    }
    cg.emit("    call *%rax")
    
    // スタッククリーンアップ
    if len(args) > 6 {
        cleanup := (len(args) - 6) * 8
        cg.emit(fmt.Sprintf("    add $%d, %%rsp", cleanup))
    }
    
    return nil
}
```

## 🎯 高度な実践課題

### 課題A-1：基本的な最適化の実装

**定数畳み込み**:
```go
func (opt *Optimizer) foldConstants(expr ast.Expression) ast.Expression {
    switch e := expr.(type) {
    case *ast.InfixExpression:
        left := opt.foldConstants(e.Left)
        right := opt.foldConstants(e.Right)
        
        if leftInt, ok := left.(*ast.IntegerLiteral); ok {
            if rightInt, ok := right.(*ast.IntegerLiteral); ok {
                switch e.Operator {
                case "+":
                    return &ast.IntegerLiteral{
                        Value: leftInt.Value + rightInt.Value,
                    }
                case "*":
                    return &ast.IntegerLiteral{
                        Value: leftInt.Value * rightInt.Value,
                    }
                }
            }
        }
    }
    return expr
}
```

### 課題A-2：簡単なガベージコレクション

**マーク・アンド・スイープ GC**:
```go
type GarbageCollector struct {
    allocatedObjects []object.Object
    reachableObjects map[object.Object]bool
}

func (gc *GarbageCollector) collect(env *object.Environment) {
    // マークフェーズ：到達可能オブジェクトをマーク
    gc.reachableObjects = make(map[object.Object]bool)
    gc.markReachable(env)
    
    // スイープフェーズ：到達不可能オブジェクトを解放
    newAllocated := []object.Object{}
    for _, obj := range gc.allocatedObjects {
        if gc.reachableObjects[obj] {
            newAllocated = append(newAllocated, obj)
        } else {
            // オブジェクトを解放
            gc.deallocate(obj)
        }
    }
    gc.allocatedObjects = newAllocated
}
```

### 課題A-3：簡単なプロファイラ

**実行時間測定**:
```go
type Profiler struct {
    functionTimes map[string]time.Duration
    callStack     []string
    startTimes    map[string]time.Time
}

func (p *Profiler) enterFunction(name string) {
    p.callStack = append(p.callStack, name)
    p.startTimes[name] = time.Now()
}

func (p *Profiler) exitFunction(name string) {
    if len(p.callStack) > 0 {
        elapsed := time.Since(p.startTimes[name])
        p.functionTimes[name] += elapsed
        p.callStack = p.callStack[:len(p.callStack)-1]
    }
}

func (p *Profiler) report() {
    fmt.Println("Function execution times:")
    for name, duration := range p.functionTimes {
        fmt.Printf("  %s: %v\n", name, duration)
    }
}
```

## 📊 性能測定課題

### ベンチマーク実装

**フィボナッチ数列による性能比較**:
```go
func BenchmarkFibonacci(b *testing.B) {
    testCases := []int{10, 20, 25, 30}
    
    for _, n := range testCases {
        input := fmt.Sprintf(`
        let fib = fn(n) {
            if (n < 2) { return n; }
            return fib(n-1) + fib(n-2);
        };
        fib(%d);
        `, n)
        
        b.Run(fmt.Sprintf("Interpreter-Fib%d", n), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                runInterpreter(input)
            }
        })
        
        b.Run(fmt.Sprintf("Compiler-Fib%d", n), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                runCompiler(input)
            }
        })
    }
}
```

## 🎮 チャレンジ課題

### チャレンジ1：簡単なREPLデバッガ

**ブレークポイントとステップ実行**:
```go
type Debugger struct {
    breakpoints map[int]bool
    currentLine int
    stepMode    bool
}

func (d *Debugger) shouldBreak(line int) bool {
    return d.breakpoints[line] || d.stepMode
}

func (d *Debugger) debugRepl(program *ast.Program, env *object.Environment) {
    scanner := bufio.NewScanner(os.Stdin)
    
    for {
        fmt.Print("(debug) ")
        if !scanner.Scan() {
            break
        }
        
        command := scanner.Text()
        switch {
        case strings.HasPrefix(command, "break "):
            // ブレークポイント設定
        case command == "step":
            d.stepMode = true
        case command == "continue":
            d.stepMode = false
        case command == "print ":
            // 変数の値を表示
        }
    }
}
```

### チャレンジ2：言語拡張

**クロージャの実装**:
```dog
let makeCounter = fn() {
    let count = 0;
    return fn() {
        count = count + 1;
        return count;
    };
};

let counter = makeCounter();
counter(); // 1
counter(); // 2
```

---

**🎯 実践課題への取り組み方**

1. **段階的に進める**: 簡単な課題から始めて徐々に複雑な課題に挑戦
2. **テストを書く**: 各機能に対してテストケースを作成
3. **性能を測定**: 実装前後で性能の変化を確認
4. **他の実装と比較**: 既存の言語処理系と比較検討
5. **文書化**: 実装の設計判断と学びを記録

**これらの課題を通じて、実践的なコンパイラ実装スキルを身につけましょう！**